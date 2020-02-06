package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewHubOptions(t *testing.T) {
	options := wspubsub.NewHubOptions()
	require.NotZero(t, options.ShutdownTimeout)
	require.False(t, options.IsDebug)
	require.NotZero(t, options.DebugFuncTimeLimit)
}
