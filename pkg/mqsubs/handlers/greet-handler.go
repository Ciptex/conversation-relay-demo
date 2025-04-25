package handlers

import (
	llmhandlers "conversation-relay/pkg/llm-handlers"
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"strings"
)

type GreetHandler struct {
	trace trace.ITracer
	repo  repo.IRepo
	LLM   llms.ILLM
}

func NewGreetHandler(repo repo.IRepo, llm llms.ILLM) *GreetHandler {
	return &GreetHandler{
		trace: trace.GetGlobalTracer(),
		repo:  repo,
		LLM:   llm,
	}
}

func (g *GreetHandler) Handle(msg types.InternalMessage, publisher chan types.MQPublish) {
	if msg.Status == "closed" {
		return
	}
	span := trace.CreateChildSpanFrom(msg.SpanRef, "GreetHandler:"+strings.ToUpper(msg.PrevEvent)+":"+msg.Track)
	defer span.Finish()
	defer trace.ExecTimer("GreetHandler::execution time", span)()

	span.SetTag("sid", msg.CallSid)
	span.SetTag("accSid", msg.AccountSid)
	span.Dev("GreetHandler::Handle", "sid", msg.CallSid, "Data", msg.Data)
	llmHandler, _ := llmhandlers.CreateLLMHadler(llmhandlers.GENERIC_HANDLER, msg.AccountSid, msg.ConfigId, &llmhandlers.HandlerContext{
		Repo:        g.repo,
		LLM:         g.LLM,
		CallSid:     msg.CallSid,
		Span:        span,
		Transcript:  []types.MessageTranscript{{Role: "system", Message: "Greet the caller pleasently"}},
		LastMessage: g.repo.GetLastMessage(msg.CallSid),
	})
	resp, err := llmHandler.Handle()
	if err != nil {
		span.Error("GreetHandler::Handle", err)
		return
	}
	msg.Data = resp
	msg.Event = types.GenericHandler

	publisher <- types.NewMQMultipPublish([]string{types.Logger}, msg, true)
	span.Debug("GreetHandler::Handle", "accSid", msg.AccountSid, "callSid", msg.CallSid, "configId", msg.ConfigId)
}
