package wspubsub_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID_Bytes(t *testing.T) {
	require.Equal(t, clientID[:], clientID.Bytes())
}

func TestUUID_String(t *testing.T) {
	require.Equal(t, "01020304-0506-0708-090a-0b0c0d0e0f10", clientID.String())
}
