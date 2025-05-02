package handlers

import (
	llmhandlers "conversation-relay/pkg/llm-handlers"
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"strings"
)

type IntentHandler struct {
	trace trace.ITracer
	repo  repo.IRepo
	LLM   llms.ILLM
}

func NewIntentHandler(repo repo.IRepo, llm llms.ILLM) *IntentHandler {
	return &IntentHandler{
		trace: trace.GetGlobalTracer(),
		repo:  repo,
		LLM:   llm,
	}
}

func (g *IntentHandler) Handle(msg types.InternalMessage, publisher chan types.MQPublish) {
	if msg.Status == "closed" {
		return
	}
	span := trace.CreateChildSpanFrom(msg.SpanRef, "IntentHandler:"+strings.ToUpper(msg.PrevEvent)+":"+msg.Track)
	defer span.Finish()
	defer trace.ExecTimer("IntentHandler::execution time", span)()

	span.SetTag("sid", msg.CallSid)
	span.SetTag("accSid", msg.AccountSid)
	span.Dev("IntentHandler::Handle", "sid", msg.CallSid, "Data", msg.Data)
	llmHandler, _ := llmhandlers.CreateLLMHadler(llmhandlers.INTENT_HANDLER, msg.AccountSid, msg.ConfigId, &llmhandlers.HandlerContext{
		Repo:        g.repo,
		LLM:         g.LLM,
		CallSid:     msg.CallSid,
		Span:        span,
		Transcript:  g.repo.GetCallContext(msg.CallSid),
		LastMessage: g.repo.GetLastMessage(msg.CallSid),
	})
	resp, err := llmHandler.Handle()
	if err != nil {
		span.Error("IntentHandler::Handle", err)
		return
	}
	msg.Data = resp
	msg.Event = types.IntentHandler
	if resp == "HUMAN_ASSISTANCE" {
		span.Debug("IntentHandler::Handle human assistance detected", "accSid", msg.AccountSid, "callSid", msg.CallSid, "configId", msg.ConfigId)
		msg.AgentTransfer = true
		// publisher <- types.NewMQMultipPublish([]string{types.Logger}, msg, true)
	}
	publisher <- types.NewMQMultipPublish([]string{types.Logger}, msg, false)
	// span.Debug("IntentHandler::Handle", "accSid", msg.AccountSid, "callSid", msg.CallSid, "configId", msg.ConfigId)
}
