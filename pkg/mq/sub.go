package mq

import (
	"bytes"
	"context"
	"conversation-relay/pkg/constants"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"strings"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"

	// register transports, do not remove
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

type MQSub struct {
	sub       mangos.Socket
	topics    string
	publisher chan types.MQPublish
	span      trace.ISpan
}

type MsgHandler interface {
	Handle(msg types.InternalMessage, publisher chan types.MQPublish)
}

func NewSubscriber(topics []string, publisher chan types.MQPublish, span trace.ISpan) *MQSub {
	child := span.AsParent("subscriber")
	defer child.Finish()

	port := constants.MQ_PUB_PORT
	sub, err := sub.NewSocket()
	if err != nil {
		child.Error("mq::NewSubscriber::error creating subscriber", err, "topics", strings.Join(topics, ","))
		return nil
	}
	if err = sub.Dial("tcp://127.0.0.1:" + port); err != nil {
		child.Error("mq::NewSubscriber::can't dial on sub socket", err)
		return nil
	}

	for _, s := range topics {
		err = sub.SetOption(mangos.OptionSubscribe, []byte(strings.TrimSpace(s)))
		if err != nil {
			span.Error("mq::NewSubscriber::unable to subscribe to", err, "topic", s)
		}
	}
	return &MQSub{
		sub:       sub,
		topics:    strings.Join(topics, ","),
		publisher: publisher,
		span:      child,
	}
}

func (s *MQSub) Listen(ctx context.Context, msgHandler MsgHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.receiveMsg()
			if err != nil {
				s.span.Error("MQSub::Listen error receiving message for topic", err, "topic", s.topics, "error", err)
				continue
			}
			go msgHandler.Handle(msg, s.publisher)
		}
	}
}

func (s *MQSub) receiveMsg() (types.InternalMessage, error) {
	msg, err := s.sub.Recv()
	if err != nil {
		return types.InternalMessage{}, err
	}
	splitted := bytes.Split(msg, []byte("|"))
	var clntMsg types.InternalMessage
	parsdMsg, err := clntMsg.ParseMsg(splitted)
	return parsdMsg, err
}
