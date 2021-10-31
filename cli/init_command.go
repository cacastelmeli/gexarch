package cli

import (
	"errors"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/processor"
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

	processor := processor.NewTemplateProcessor(&config.ProcessorConfig{
		CliConfig: &config.CliConfig{
			TypesPath: typesPath,
		},
	})
	processor.ProcessInit()

	return nil
}
