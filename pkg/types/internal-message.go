package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type InternalMessage struct {
	Id                        string                 `json:"id"`
	Event                     string                 `json:"event"` // what event it is like transcript or summary etc
	PrevEvent                 string                 `json:"prevEvent"`
	ConfigId                  string                 `json:"configId"`
	AccountSid                string                 `json:"accountSid"`
	CallSid                   string                 `json:"callSid"`
	FromLang                  string                 `json:"fromLanguage"`
	ToLang                    string                 `json:"toLanguage"`  // once the conversion is done the conversion module set the language
	Track                     string                 `json:"track"`       // inbound or outbound etc
	Status                    string                 `json:"status"`      // can be inprogress etc. Once the event is done set the status to closed
	Data                      string                 `json:"data"`        // data in string format, can be json, plain text etc
	SayLanguage               string                 `json:"sayLanguage"` // language used to say translated message back to caller
	SayVoice                  string                 `json:"sayVoice"`
	From                      string                 `json:"from"`
	Attributes                map[string]interface{} `json:"attributes"`
	SpanRef                   map[string]string      `json:"spanRef"`
	SessionId                 string                 `json:"sessionId"`
	OrgTwilioCRInboundPayload TwilioCRInboundPayload `json:"orgTwilioCRInboundPayload"`
	IsLastMessage             bool                   `json:"isLastMessage"`
}

func NewInternalMessage(callSid string, data string) InternalMessage {
	uuid := uuid.New().String()
	return InternalMessage{
		Id:      uuid,
		CallSid: callSid,
		Data:    data,
	}
}

func (c InternalMessage) ParseMsg(msg [][]byte) (InternalMessage, error) {
	topic := string(msg[0])
	var clntMsg InternalMessage
	err := json.Unmarshal(msg[1], &clntMsg)
	clntMsg.PrevEvent = clntMsg.Event
	clntMsg.Event = topic
	return clntMsg, err
}
