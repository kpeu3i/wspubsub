package wspubsub

// MessageType enumerates possible message types.
type MessageType byte

const (
	MessageTypeText   MessageType = 1
	MessageTypeBinary MessageType = 2
	MessageTypePing   MessageType = 9
)

// Message represents a data type to send over a WebSocket connection.
type Message struct {
	Type    MessageType
	Payload []byte
}

// NewTextMessage initializes a new text Message from bytes.
func NewTextMessage(payload []byte) Message {
	return Message{Type: MessageTypeText, Payload: payload}
}

// NewTextMessage initializes a new text Message from string.
func NewTextMessageFromString(payload string) Message {
	return NewTextMessage([]byte(payload))
}

// NewBinaryMessage initializes a new binary Message from bytes.
func NewBinaryMessage(payload []byte) Message {
	return Message{Type: MessageTypeBinary, Payload: payload}
}

// NewBinaryMessageFromString initializes a new binary Message from string.
func NewBinaryMessageFromString(payload string) Message {
	return NewBinaryMessage([]byte(payload))
}

// NewPingMessage initializes a new ping Message.
func NewPingMessage() Message {
	return Message{Type: MessageTypePing}
}
