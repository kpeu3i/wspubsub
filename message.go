package wspubsub

type MessageType byte

const (
	MessageTypeText   MessageType = 1
	MessageTypeBinary MessageType = 2
	MessageTypePing   MessageType = 9
)

type Message struct {
	Type    MessageType
	Payload []byte
}

func NewTextMessage(payload []byte) Message {
	return Message{Type: MessageTypeText, Payload: payload}
}

func NewTextMessageFromString(payload string) Message {
	return NewTextMessage([]byte(payload))
}

func NewBinaryMessage(payload []byte) Message {
	return Message{Type: MessageTypeBinary, Payload: payload}
}

func NewBinaryMessageFromString(payload string) Message {
	return NewBinaryMessage([]byte(payload))
}

func NewPingMessage() Message {
	return Message{Type: MessageTypePing}
}
