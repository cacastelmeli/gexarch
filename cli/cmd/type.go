package cmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/processor"
	"github.com/aeroxmotion/gexarch/util"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	routerExprSource   = "%sRouter.New%sRouter(conf)"
	subRoutersFilename = "%s/shared/infrastructure/routers/sub_routers.go"
	subRouterPath      = "%s/%s/infrastructure/router"
)

func TypeCommand() *cli.Command {
	return &cli.Command{
		Name:    "type",
		Aliases: []string{"t"},
		Usage:   "Generate scaffold by `type`",
		Action:  typeCommandAction,
	}
}

func typeCommandAction(ctx *cli.Context) error {
	targetType := strcase.ToCamel(ctx.Args().Get(0))

	if targetType == "" {
		return errors.New("missing `type` name")
	}

	conf := config.GetProcessorConfigByType(targetType)
	templateProcessor, codemodProcessor := processor.NewTemplateProcessor(conf), processor.NewCodemodProcessor(conf)

	templateProcessor.ProcessTemplate(
		"type",
		path.Join(conf.TypesPath, strcase.ToSnake(conf.TypeName)),
	)
	codemodProcessor.ProcessFile(fmt.Sprintf(subRoutersFilename, conf.TypesPath), addRouterInstance)

	return nil
}

func addRouterInstance(fileSet *token.FileSet, fileNode *ast.File, conf *config.ProcessorConfig) ast.Node {
	// Add top-level named import
	// Will compile into something like:
	// import <TypeName>Router "<ModulePath>/<TypeName>/..."
	astutil.AddNamedImport(
		fileSet,
		fileNode,
		fmt.Sprintf("%sRouter", strcase.ToLowerCamel(conf.TypeName)),
		fmt.Sprintf(subRouterPath, conf.ModulePath, strcase.ToSnake(conf.TypeName)),
	)

	// Attach router's instance call expression
	// to the list of registered routers:
	//
	// []rest.SharedRouter{
	//     <TypeName>Router.New<TypeName>Router(conf),
	// }
	postTransform := func(c *astutil.Cursor) bool {
		compositeLit, inCompositeLit := c.Parent().(*ast.CompositeLit)

		if !inCompositeLit {
			return true
		}

		arrayType, inArrayType := compositeLit.Type.(*ast.ArrayType)

		if !inArrayType {
			return true
		}

		selectorExpr, hasSelectorExpr := arrayType.Elt.(*ast.SelectorExpr)

		if !hasSelectorExpr || selectorExpr.Sel.Name != "SharedRouter" {
			return true
		}

		routerExpr, err := parser.ParseExpr(
			fmt.Sprintf(
				routerExprSource,
				strcase.ToLowerCamel(conf.TypeName),
				strcase.ToCamel(conf.TypeName),
			),
		)
		util.PanicIfError(err)

		compositeLit.Elts = append(compositeLit.Elts, routerExpr)
		return false
	}

	return astutil.Apply(fileNode, func(c *astutil.Cursor) bool {
		return true
	}, postTransform)
}
