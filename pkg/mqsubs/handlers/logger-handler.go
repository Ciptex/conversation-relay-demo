package handlers

import (
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"strings"
)

type MsgLogger struct {
	trace trace.ITracer
}

func NewLoggerHandler() *MsgLogger {
	return &MsgLogger{
		trace: trace.GetGlobalTracer(),
	}
}

func (l *MsgLogger) Handle(msg types.InternalMessage, publisher chan types.MQPublish) {
	if msg.Status == "closed" {
		return
	}
	span := trace.CreateChildSpanFrom(msg.SpanRef, "MsgLogger:"+strings.ToUpper(msg.PrevEvent)+":"+msg.Track)
	defer span.Finish()

	span.SetTag("sid", msg.CallSid)
	span.SetTag("accSid", msg.AccountSid)
	span.Dev("LoggerHandler", "accSid", msg.AccountSid, "sid", msg.CallSid, "Data", msg.Data)
}
