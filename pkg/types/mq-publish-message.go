package types

type MQPublish struct {
	Topics         []string
	Message        InternalMessage
	BroadcastToHub bool
}

func NewMQPublish(topic string, clientMsg InternalMessage, broadcastToHub bool) MQPublish {
	return MQPublish{
		Topics:         []string{topic},
		Message:        clientMsg,
		BroadcastToHub: broadcastToHub,
	}
}

func NewMQMultipPublish(topics []string, clientMsg InternalMessage, broadcastToClient bool) MQPublish {
	return MQPublish{
		Topics:         topics,
		Message:        clientMsg,
		BroadcastToHub: broadcastToClient,
	}
}
