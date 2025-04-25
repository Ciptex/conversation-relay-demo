package llmhandlers

import (
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"

	"errors"
)

type LLMHandler interface {
	Handle() (string, error)
}

type HandlerType string

const (
	GENERIC_HANDLER HandlerType = "GENERIC_HANDLER"
	GREET_HANDLER   HandlerType = "GREET_HANDLER"
)

type HandlerContext struct {
	Repo        repo.IRepo
	LLM         llms.ILLM
	CallSid     string
	Span        trace.ISpan
	Transcript  []types.MessageTranscript
	LastMessage string
}

func CreateLLMHadler(handlerType HandlerType, accSid, configSid string, context *HandlerContext) (LLMHandler, error) {
	switch handlerType {
	case GENERIC_HANDLER:
		return newGenericPromptHandler(accSid, configSid, context), nil
	case GREET_HANDLER:
		return newGreetPromptHandler(accSid, configSid, context), nil
	default:
		return nil, errors.New("unable to find handler for " + string(handlerType))
	}
}
