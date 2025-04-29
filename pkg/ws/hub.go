package ws

import (
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/types"
)

type Hub struct {
	twiliClients  map[string]*TwilioClient
	VoiceResponse chan types.InternalMessage
	publisher     chan types.MQPublish
	repo          repo.IRepo
}

func NewHub(repo repo.IRepo) *Hub {
	return &Hub{
		twiliClients:  make(map[string]*TwilioClient),
		VoiceResponse: make(chan types.InternalMessage),
		repo:          repo,
	}
}

func (h *Hub) SetMQPublisher(publisher chan types.MQPublish) {
	h.publisher = publisher
}

func (h *Hub) Publish(topics []string, msg types.InternalMessage, broadcastToClient bool) {
	if h.publisher == nil {
		return
	}
	h.publisher <- types.NewMQMultipPublish(topics, msg, broadcastToClient)
}

func (h *Hub) DB() repo.IRepo {
	return h.repo
}

func (h *Hub) Listen() {
	for {
		select {
		case msg := <-h.VoiceResponse:
			if client, ok := h.twiliClients[msg.CallSid]; ok {
				client.respond(msg)
			}
		}
	}
}

func (h *Hub) addTwilioClient(client *TwilioClient) {
	h.twiliClients[client.callSid] = client
}

func (h *Hub) removeTwilioClient(client *TwilioClient) {
	delete(h.twiliClients, client.callSid)
}
