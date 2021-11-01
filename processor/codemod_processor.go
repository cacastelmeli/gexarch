package processor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/util"
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/ast/astutil"
)

type CodemodProcessor struct {
}

const (
	routerExprSource   = "%sRouter.New%sRouter(conf)"
	subRoutersFilename = "%s/shared/infrastructure/routers/sub_routers.go"
	subRouterPath      = "%s/%s/infrastructure/router"
)

func NewCodemodProcessor() *CodemodProcessor {
	return &CodemodProcessor{}
}

func (processor *CodemodProcessor) ProcessType(conf *config.ProcessorConfig) {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, fmt.Sprintf(subRoutersFilename, conf.TypesPath), nil, parser.ParseComments)
	util.PanicIfError(err)

	astutil.AddNamedImport(
		fileSet,
		node,
		fmt.Sprintf("%sRouter", strcase.ToLowerCamel(conf.TypeName)),
		fmt.Sprintf(subRouterPath, conf.ModulePath, strcase.ToSnake(conf.TypeName)),
	)

	result := astutil.Apply(node, func(c *astutil.Cursor) bool {
		return true
	}, func(c *astutil.Cursor) bool {
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
	})

	f, err := os.Create(fmt.Sprintf(subRoutersFilename, conf.TypesPath))
	util.PanicIfError(err)
	defer f.Close()

	util.PanicIfError(printer.Fprint(f, fileSet, result))
}
