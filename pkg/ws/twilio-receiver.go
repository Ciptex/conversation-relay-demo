package ws

import (
	"bytes"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type TwilioClient struct {
	conn      *websocket.Conn
	callSid   string
	accSid    string
	configSid string
	span      trace.ISpan
	hub       *Hub
}

func newTwilioClient(hub *Hub, conn *websocket.Conn, configSid string) *TwilioClient {
	return &TwilioClient{
		conn:      conn,
		configSid: configSid,
		hub:       hub,
	}
}

func (c *TwilioClient) listen() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.span.Error("TwilioStream::listen::error reading message", err)
			break
		}
		go c.processMessage(message)
	}
	c.span.Info("TwilioStream::listen twilio call ended", "accSid", c.accSid, "callSid", c.callSid)
	c.span.Finish()
	c.hub.removeTwilioClient(c)
}

func (c *TwilioClient) parseCustomParameters(payload types.TwilioCRInboundPayload) {
	if payload.CustomParameters == nil {
		return
	}
	for key, value := range payload.CustomParameters {
		switch key {
		case "accSid":
			c.accSid = value
		case "configSid":
			c.configSid = value
		}
	}
}

func (c *TwilioClient) processMessage(message []byte) {
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	var payload types.TwilioCRInboundPayload
	err := json.Unmarshal(message, &payload)
	if err != nil {
		c.span.Error("TwilioStream::processMessage::error unmarshalling message", err)
		return
	}
	if payload.Type == "setup" {
		c.accSid = payload.AccountSid
		c.callSid = payload.CallSid
		c.parseCustomParameters(payload)
		c.hub.addTwilioClient(c)
		c.hub.repo.GetAccountConfig(c.accSid, c.configSid)
		go c.publish(payload, []string{types.GreetHandler})
		c.span.Info("TwilioStream::processMessage::setup", "accSid", c.accSid, "configId", c.configSid, "callSid", payload.CallSid)
	} else if payload.Type == "prompt" {
		c.hub.DB().AddCallContext(c.callSid, "human", payload.VoicePrompt)
		c.span.Info("TwilioStream::processMessage::prompt", "message", payload.VoicePrompt, "lang", payload.Lang)
		go c.publish(payload, []string{types.Logger, types.GenericHandler})
	}
}

func (c *TwilioClient) publish(msg types.TwilioCRInboundPayload, subs []string) {
	m := types.NewInternalMessage(c.callSid, strings.TrimSpace(msg.VoicePrompt))
	m.AccountSid = c.accSid
	m.ConfigId = c.configSid
	m.FromLang = msg.Lang
	m.ToLang = msg.Lang
	m.SessionId = msg.SessionId
	m.OrgTwilioCRInboundPayload = msg
	m.Track = msg.Direction
	m.IsLastMessage = true
	m.From = msg.From
	c.hub.Publish(subs, m, false)
}

func (c *TwilioClient) respond(msg types.InternalMessage) {
	reply := map[string]interface{}{
		// "type": "end",
		// "handoffData": "{\"reasonCode\":\"live-agent-handoff\", \"reason\": \"The caller wants to talk to a real person\"}",
		"type":  "text",
		"token": msg.Data,
		"last":  msg.IsLastMessage,
	}
	replyStr, _ := json.Marshal(reply)
	c.hub.DB().AddCallContext(c.callSid, "ai", msg.Data)
	c.conn.WriteMessage(websocket.TextMessage, replyStr)
	c.span.Info("TwilioStream::Respond", "accSid", c.accSid, "callSid", c.callSid, "message", msg.Data)
}

func StartTwilioHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	span := trace.GetGlobalTracer().Start("TwilioConversationRelay")
	defer span.Finish()
	r.ParseForm()
	vars := mux.Vars(r)
	configId := vars["configId"]
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		span.Error("TwilioStream::StartTwilioHandler::error upgrading request", err)
		return
	}
	client := newTwilioClient(hub, conn, configId)
	client.span = span
	go client.listen()
}
