package processor

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/util"
)

type TransformFunc func(fileSet *token.FileSet, file *ast.File, conf *config.ProcessorConfig) ast.Node

type CodemodProcessor struct {
	fileSet *token.FileSet
	config  *config.ProcessorConfig
}

func NewCodemodProcessor(conf *config.ProcessorConfig) *CodemodProcessor {
	return &CodemodProcessor{
		fileSet: token.NewFileSet(),
		config:  conf,
	}
}

func (processor *CodemodProcessor) ProcessFile(filename string, transformFunc TransformFunc) {
	// Parse given `filename` into an AST
	// TODO: We can optimize parsing by passing a bit flag (IDK)
	node, err := parser.ParseFile(processor.fileSet, filename, nil, parser.ParseComments)
	util.PanicIfError(err)

	f, err := os.Create(filename)
	util.PanicIfError(err)
	defer f.Close()

	util.PanicIfError(
		printer.Fprint(
			f,
			processor.fileSet,
			transformFunc(processor.fileSet, node, processor.config),
		),
	)
}
