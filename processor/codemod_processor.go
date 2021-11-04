package processor

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/util"
)

const (
	// Skip object resolution for fast path
	optimizedParserMode parser.Mode = parser.ParseComments | parser.SkipObjectResolution
)

type TransformFunc func(fileSet *token.FileSet, fileNode *ast.File, conf *config.ProcessorConfig) ast.Node

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
	fileNode, err := parser.ParseFile(processor.fileSet, filename, nil, optimizedParserMode)
	util.PanicIfError(err)

	file, err := os.Create(filename)
	util.PanicIfError(err)
	defer file.Close()

	err = format.Node(
		file,
		processor.fileSet,
		transformFunc(processor.fileSet, fileNode, processor.config),
	)
	util.PanicIfError(err)
}
