package llms

import (
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
)

type ILLM interface {
	New(context types.LLMModelContext) LLMModel
}

type LLMModel interface {
	CreateChatCompletion(config types.AccountConfig, sid, prompt string, span trace.ISpan) (string, error)
	CreateEmbedding(config types.AccountConfig, text string, span trace.ISpan) ([]float32, error)
}

func CreateLLMModel() ILLM {
	return newGTPLLM()
}
