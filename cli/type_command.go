package cli

import (
	"errors"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/processor"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

func typeCommand() *cli.Command {
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

	tplProcessor := processor.NewTemplateProcessor(conf)
	tplProcessor.ProcessByType()

	codemodProcessor := processor.NewCodemodProcessor()
	codemodProcessor.ProcessType(conf)

	return nil
}
