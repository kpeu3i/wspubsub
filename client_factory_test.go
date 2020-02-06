package wspubsub_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kpeu3i/wspubsub"
	"github.com/kpeu3i/wspubsub/mock"
	"github.com/stretchr/testify/require"
)

func TestClientFactory_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGenerator := mock.NewMockUUIDGenerator(ctrl)
	uuidGenerator.
		EXPECT().
		GenerateV4().
		Return(clientID)

	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)
	clientOptions := wspubsub.NewClientOptions()
	logger := mock.NewMockLogger(ctrl)
	factory := wspubsub.NewClientFactory(clientOptions, uuidGenerator, upgrader, logger)
	client := factory.Create()

	require.NotNil(t, client)
	require.NotNil(t, client.ID())
}
