package processor

import "github.com/aeroxmotion/gexarch/config"

type HookProcessFunc func(templateProcessor *TemplateProcessor, codemodProcessor *CodemodProcessor)

func Process(conf *config.ProcessorConfig, hookFunc HookProcessFunc) {
	hookFunc(
		NewTemplateProcessor(conf),
		NewCodemodProcessor(conf),
	)
}
