package mq

import (
	"bytes"
	"conversation-relay/pkg/constants"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"encoding/json"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
)

type MQPub struct {
	pub           mangos.Socket
	Publisher     chan types.MQPublish
	voiceResponse chan types.InternalMessage
	toHub         chan types.InternalMessage
}

func NewPublisher(voiceResponse chan types.InternalMessage, span trace.ISpan) *MQPub {
	child := span.AsParent("publisher")
	defer child.Finish()

	child.Info("mq::start publisher")
	pub, err := pub.NewSocket()
	if err != nil {
		child.Error("mq::NewPublisher::error creating mq publisher", err)
		return nil
	}
	if err = pub.Listen("tcp://127.0.0.1:" + constants.MQ_PUB_PORT); err != nil {
		child.Error("mq::NewPublisher::can't listen on pub socket:", err)
	}
	child.Info("mq::publisher started successfully")
	return &MQPub{
		pub:           pub,
		Publisher:     make(chan types.MQPublish, 1000),
		toHub:         make(chan types.InternalMessage, 1000),
		voiceResponse: voiceResponse,
	}
}

func (mq *MQPub) Listen() {
	for {
		select {
		case msg := <-mq.Publisher:
			for _, topic := range msg.Topics {
				mq.publish(topic, msg.Message)
			}
			if msg.BroadcastToHub {
				mq.voiceResponse <- msg.Message
			}
		}
	}
}

func (mq *MQPub) publish(topic string, message types.InternalMessage) {
	json, _ := json.Marshal(message)
	var msgArr [][]byte
	msgArr = append(msgArr, []byte(topic))
	msgArr = append(msgArr, json)
	joinedBytes := bytes.Join(msgArr, []byte("|"))
	mq.pub.Send(joinedBytes)
}
