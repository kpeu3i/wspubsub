package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewLogrusOptions(t *testing.T) {
	options := wspubsub.NewLogrusOptions()
	require.NotZero(t, options.Output)
}
