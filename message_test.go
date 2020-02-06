package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewTextMessage(t *testing.T) {
	b := []byte("TEST")
	message := wspubsub.NewTextMessage(b)
	require.Equal(t, wspubsub.MessageTypeText, message.Type)
	require.Equal(t, b, message.Payload)

	s := "TEST"
	message = wspubsub.NewTextMessageFromString(s)
	require.Equal(t, wspubsub.MessageTypeText, message.Type)
	require.Equal(t, []byte(s), message.Payload)
}

func TestNewBinaryMessage(t *testing.T) {
	b := []byte("TEST")
	message := wspubsub.NewBinaryMessage(b)
	require.Equal(t, wspubsub.MessageTypeBinary, message.Type)
	require.Equal(t, b, message.Payload)

	s := "TEST"
	message = wspubsub.NewBinaryMessageFromString(s)
	require.Equal(t, wspubsub.MessageTypeBinary, message.Type)
	require.Equal(t, []byte(s), message.Payload)
}
