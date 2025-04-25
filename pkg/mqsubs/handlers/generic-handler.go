package handlers

import (
	llmhandlers "conversation-relay/pkg/llm-handlers"
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"strings"
)

type GenericHandler struct {
	trace trace.ITracer
	repo  repo.IRepo
	LLM   llms.ILLM
}

func NewGenericHandler(repo repo.IRepo, llm llms.ILLM) *GenericHandler {
	return &GenericHandler{
		trace: trace.GetGlobalTracer(),
		repo:  repo,
		LLM:   llm,
	}
}

func (g *GenericHandler) Handle(msg types.InternalMessage, publisher chan types.MQPublish) {
	if msg.Status == "closed" {
		return
	}
	span := trace.CreateChildSpanFrom(msg.SpanRef, "GenericHandler:"+strings.ToUpper(msg.PrevEvent)+":"+msg.Track)
	defer span.Finish()
	defer trace.ExecTimer("GenericHandler::execution time", span)()

	span.SetTag("sid", msg.CallSid)
	span.SetTag("accSid", msg.AccountSid)
	span.Dev("GenericHandler::Handle", "sid", msg.CallSid, "Data", msg.Data)
	llmHandler, _ := llmhandlers.CreateLLMHadler(llmhandlers.GENERIC_HANDLER, msg.AccountSid, msg.ConfigId, &llmhandlers.HandlerContext{
		Repo:        g.repo,
		LLM:         g.LLM,
		CallSid:     msg.CallSid,
		Span:        span,
		Transcript:  g.repo.GetCallContext(msg.CallSid),
		LastMessage: g.repo.GetLastMessage(msg.CallSid),
	})
	resp, err := llmHandler.Handle()
	if err != nil {
		span.Error("GenericHandler::Handle", err)
		return
	}
	msg.Data = resp
	msg.Event = types.GenericHandler

	publisher <- types.NewMQMultipPublish([]string{types.Logger}, msg, true)
	span.Debug("GenericHandler::Handle", "accSid", msg.AccountSid, "callSid", msg.CallSid, "configId", msg.ConfigId)
}
