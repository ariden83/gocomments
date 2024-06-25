package comments

import (
	"go/ast"
)

type commentsProcess interface {
	isActive() bool
	commentConst(string, bool) (string, error)
	commentFunc(fn *ast.FuncDecl) (string, error)
	commentType(genDecl *ast.GenDecl) (string, error)
	commentVar(name, declType, explainVar string, exported bool) (string, error)
}

func newProcessor(cfg *CommentConfig) commentsProcess {
	if cfg == nil {
		return &defaultProcess{
			activeExamples: cfg.ActiveExamples,
		}
	}

	openAIProcess := openAI{
		OpenAIConfig: cfg.OpenAI,
	}
	if openAIProcess.isActive() {
		return &openAIProcess
	}

	anthropicProcess := anthropic{
		AnthropicConfig: cfg.Anthropic,
	}
	if anthropicProcess.isActive() {
		return &anthropicProcess
	}

	return &defaultProcess{}
}
