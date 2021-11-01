package cli

import (
	"errors"
	"os"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/processor"
	"github.com/aeroxmotion/gexarch/util"
	"github.com/urfave/cli/v2"
)

func initCommand() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize gexarch (this will override previous configuration)",
		Action:  initCommandAction,
	}
}

func initCommandAction(ctx *cli.Context) error {
	typesPath := ctx.Args().Get(0)

	if typesPath == "" {
		return errors.New("missing types-path argument")
	}

	conf := &config.ProcessorConfig{
		CliConfig: &config.CliConfig{
			TypesPath: typesPath,
		},
		ModulePath: util.ParseModfile().Module.Mod.Path,
	}

	processor.Process(conf, func(templateProcessor *processor.TemplateProcessor, codemodProcessor *processor.CodemodProcessor) {
		workingDirectory, err := os.Getwd()
		util.PanicIfError(err)

		templateProcessor.ProcessTemplate("init", workingDirectory)
	})

	return nil
}
