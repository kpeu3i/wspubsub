package wspubsub_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kpeu3i/wspubsub"
	"github.com/kpeu3i/wspubsub/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestClientStore_Clients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	clientsNum := 10
	clients := make(map[wspubsub.UUID]*wspubsub.Client, clientsNum)
	clientList := make([]*wspubsub.Client, 0, clientsNum)
	for i := 0; i < clientsNum; i++ {
		cid := wspubsub.UUID([16]byte{byte(i), 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
		clientOptions := wspubsub.NewClientOptions()
		client := wspubsub.NewClient(clientOptions, cid, upgrader, logger)
		clients[cid] = client
		clientList = append(clientList, client)
	}

	clientStoreOptions := wspubsub.NewClientStoreOptions()
	clientStore := wspubsub.NewClientStore(clientStoreOptions, logger)

	t.Run("Check default state of a storage", func(t *testing.T) {
		client, err := clientStore.Get(clientID)
		require.Nil(t, client)
		require.Equal(t, wspubsub.NewClientNotFoundError(clientID), errors.Cause(err).(*wspubsub.ClientNotFoundError))
		require.Equal(t, 0, clientStore.Count())

		err = clientStore.Unset(clientID)
		require.NoError(t, err)
	})

	t.Run("Fill a storage", func(t *testing.T) {
		i := 1
		for cid, client := range clients {
			clientStore.Set(client)
			storedClient, err := clientStore.Get(cid)
			require.NoError(t, err)
			require.Equal(t, client, storedClient)
			require.Equal(t, i, clientStore.Count())

			i++
		}
	})

	t.Run("Remove a client from a storage", func(t *testing.T) {
		removedClientID := clientList[0].ID()
		err := clientStore.Unset(removedClientID)
		require.NoError(t, err)
		require.Equal(t, len(clientList)-1, clientStore.Count())

		client, err := clientStore.Get(removedClientID)
		require.Nil(t, client)
		require.Equal(t, wspubsub.NewClientNotFoundError(removedClientID), errors.Cause(err).(*wspubsub.ClientNotFoundError))

		for cid, client := range clients {
			if client.ID() == removedClientID {
				continue
			}

			storedClient, err := clientStore.Get(cid)
			require.NoError(t, err)
			require.Equal(t, client, storedClient)
		}
	})
}

func TestClientStore_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	clientStoreOptions := wspubsub.NewClientStoreOptions()
	clientStore := wspubsub.NewClientStore(clientStoreOptions, logger)

	numClients := 100
	numClientsInChannels := 0
	availableChannels := []string{"X", "Y", "Z"}
	clients := map[string]map[wspubsub.UUID]*wspubsub.Client{}

	for i := 0; i < numClients; i++ {
		cid := wspubsub.UUID([16]byte{byte(i), 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
		clientOptions := wspubsub.NewClientOptions()
		client := wspubsub.NewClient(clientOptions, cid, upgrader, logger)
		clientStore.Set(client)

		numChannels := rand.Intn(len(availableChannels)-1) + 1
		channels := make([]string, 0, numChannels)
		for j := 0; j < numChannels; j++ {
			channels = append(channels, availableChannels[rand.Intn(len(availableChannels))])
		}

		err := clientStore.SetChannels(cid, channels...)
		require.NoError(t, err)

		for _, channel := range channels {
			if _, ok := clients[channel]; !ok {
				clients[channel] = make(map[wspubsub.UUID]*wspubsub.Client)
			}

			if _, ok := clients[channel][client.ID()]; !ok {
				numClientsInChannels++
			}

			clients[channel][client.ID()] = client
		}
	}

	t.Run("Count clients in a channel", func(t *testing.T) {
		for _, channel := range availableChannels {
			require.Equal(t, len(clients[channel]), clientStore.Count(channel))
		}
	})

	t.Run("Count clients in a few random channels", func(t *testing.T) {
		channel1 := availableChannels[rand.Intn(len(availableChannels))]
		channel2 := availableChannels[rand.Intn(len(availableChannels))]
		require.Equal(t, len(clients[channel1])+len(clients[channel2]), clientStore.Count(channel1, channel2))
	})

	t.Run("Count clients in all channels", func(t *testing.T) {
		require.Equal(t, numClientsInChannels, clientStore.Count(availableChannels...))
	})

	t.Run("Count clients without channels", func(t *testing.T) {
		require.Equal(t, numClients, clientStore.Count())
	})

	t.Run("Find clients in a random channel", func(t *testing.T) {
		foundClients := 0
		channel := availableChannels[rand.Intn(len(availableChannels))]

		fn := func(client wspubsub.WebsocketClient) error {
			foundClients++
			_, ok := clients[channel][client.ID()]
			require.True(t, ok)

			return nil
		}

		err := clientStore.Find(fn, channel)
		require.NoError(t, err)

		require.Equal(t, len(clients[channel]), foundClients)
	})

	t.Run("Find clients in a few random channels", func(t *testing.T) {
		numFoundClients := 0
		channel1 := availableChannels[rand.Intn(len(availableChannels))]
		channel2 := availableChannels[rand.Intn(len(availableChannels))]

		fn := func(client wspubsub.WebsocketClient) error {
			numFoundClients++
			_, ok1 := clients[channel1][client.ID()]
			_, ok2 := clients[channel2][client.ID()]

			require.True(t, ok1 || ok2)

			return nil
		}

		err := clientStore.Find(fn, channel1, channel2)
		require.NoError(t, err)

		require.Equal(t, len(clients[channel1])+len(clients[channel2]), numFoundClients)
	})

	t.Run("Find clients in all channels", func(t *testing.T) {
		numFoundClients := 0
		fn := func(client wspubsub.WebsocketClient) error {
			numFoundClients++
			exists := false
			for _, channel := range availableChannels {
				if _, ok := clients[channel][client.ID()]; ok {
					exists = true
					break
				}
			}

			require.True(t, exists)

			return nil
		}

		err := clientStore.Find(fn, availableChannels...)
		require.NoError(t, err)

		require.Equal(t, numClientsInChannels, numFoundClients)
	})

	t.Run("Find clients without channels", func(t *testing.T) {
		foundClients := 0
		fn := func(client wspubsub.WebsocketClient) error {
			foundClients++

			exists := false
			for _, channel := range availableChannels {
				if _, ok := clients[channel][client.ID()]; ok {
					exists = true
					break
				}
			}

			require.True(t, exists)

			return nil
		}

		err := clientStore.Find(fn)
		require.NoError(t, err)

		require.Equal(t, numClients, foundClients)
	})

	t.Run("Find clients with iteration error", func(t *testing.T) {
		fn := func(client wspubsub.WebsocketClient) error {
			return errors.New("something went wrong")
		}

		err := clientStore.Find(fn)
		require.Error(t, err)
	})
}

func TestClientStore_Channels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	clientStoreOptions := wspubsub.NewClientStoreOptions()
	clientStore := wspubsub.NewClientStore(clientStoreOptions, logger)

	clientOptions := wspubsub.NewClientOptions()
	client := wspubsub.NewClient(clientOptions, clientID, upgrader, logger)
	clientStore.Set(client)

	unknownClientID := wspubsub.UUID{}

	t.Run("Use channels for unknown client", func(t *testing.T) {
		channels, err := clientStore.Channels(unknownClientID)
		require.Empty(t, channels)
		require.Equal(t, wspubsub.NewClientNotFoundError(unknownClientID), errors.Cause(err).(*wspubsub.ClientNotFoundError))

		numChannels, err := clientStore.CountChannels(unknownClientID)
		require.Equal(t, 0, numChannels)
		require.Equal(t, wspubsub.NewClientNotFoundError(unknownClientID), errors.Cause(err).(*wspubsub.ClientNotFoundError))

		err = clientStore.SetChannels(unknownClientID, "X")
		require.Equal(t, wspubsub.NewClientNotFoundError(unknownClientID), errors.Cause(err).(*wspubsub.ClientNotFoundError))

		err = clientStore.UnsetChannels(unknownClientID, "X")
		require.Equal(t, wspubsub.NewClientNotFoundError(unknownClientID), errors.Cause(err).(*wspubsub.ClientNotFoundError))
	})

	t.Run("Check default sate of channels", func(t *testing.T) {
		channels, err := clientStore.Channels(clientID)
		require.NoError(t, err)
		require.Empty(t, channels)
		numChannels, err := clientStore.CountChannels(clientID)
		require.NoError(t, err)
		require.Equal(t, 0, numChannels)
	})

	t.Run("Set empty channels", func(t *testing.T) {
		err := clientStore.SetChannels(clientID)
		require.NoError(t, err)
	})

	t.Run("Set channels", func(t *testing.T) {
		channels := []string{"X", "Y", "W", "Z"}
		sort.Strings(channels)

		err := clientStore.SetChannels(clientID, channels...)
		require.NoError(t, err)
		newChannels, err := clientStore.Channels(clientID)
		sort.Strings(newChannels)
		require.NoError(t, err)
		require.Equal(t, channels, newChannels)
		numChannels, err := clientStore.CountChannels(clientID)
		require.NoError(t, err)
		require.Equal(t, len(channels), numChannels)
	})

	t.Run("Unset channels", func(t *testing.T) {
		channels := []string{"X", "Z"}
		sort.Strings(channels)

		err := clientStore.UnsetChannels(clientID, "Y", "W")
		require.NoError(t, err)
		newChannels, err := clientStore.Channels(clientID)
		sort.Strings(newChannels)
		require.NoError(t, err)
		require.Equal(t, channels, newChannels)
		numChannels, err := clientStore.CountChannels(clientID)
		require.NoError(t, err)
		require.Equal(t, len(channels), numChannels)

		err = clientStore.UnsetChannels(clientID)
		require.NoError(t, err)
		newChannels, err = clientStore.Channels(clientID)
		sort.Strings(newChannels)
		require.NoError(t, err)
		require.Empty(t, newChannels)
	})
}
