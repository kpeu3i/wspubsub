package wspubsub

import (
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/pkg/errors"
)

type IterateFunc func(client WebsocketClient) error

type clientsBuffer struct {
	clients []WebsocketClient
}

type ClientStore struct {
	options           ClientStoreOptions
	logger            Logger
	clientsShardList  []*clientStoreClientsShard
	channelsShardList []*clientStoreChannelsShard
	clientsPool       sync.Pool
}

func (s *ClientStore) Get(clientID UUID) (WebsocketClient, error) {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.get: took=%s", end)
			}
		}()
	}

	clientsShard := s.clientsShard(clientID)

	client, err := clientsShard.Get(clientID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client, nil
}

func (s *ClientStore) Set(client WebsocketClient) {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.set: took=%s", end)
			}
		}()
	}

	clientsShard := s.clientsShard(client.ID())
	clientsShard.Set(client)
}

func (s *ClientStore) Unset(clientID UUID) error {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.unset: took=%s", end)
			}
		}()
	}

	clientsShard := s.clientsShard(clientID)
	clientsShard.Unset(clientID)

	for _, channelsShard := range s.channelsShardList {
		channelsShard.Unlink(clientID)
	}

	return nil
}

func (s *ClientStore) Count(channels ...string) int {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.count: took=%s", end)
			}
		}()
	}

	count := 0
	if len(channels) == 0 {
		for _, clientsShard := range s.clientsShardList {
			count += clientsShard.Count()
		}

		return count
	}

	for _, channel := range channels {
		channelsShard := s.channelsShard(channel)
		count += channelsShard.Count(channel)
	}

	return count
}

func (s *ClientStore) Find(fn IterateFunc, channels ...string) error {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.find: took=%s", end)
			}
		}()
	}

	buff := s.clientsPool.Get().(*clientsBuffer)

	if len(channels) == 0 {
		for _, clientsShard := range s.clientsShardList {
			clientsShard.Iterate(func(client WebsocketClient) {
				buff.clients = append(buff.clients, client)
			})
		}
	} else {
		for _, channel := range channels {
			channelsShard := s.channelsShard(channel)
			channelsShard.Iterate(func(client WebsocketClient, channel string) {
				buff.clients = append(buff.clients, client)
			})
		}
	}

	for _, client := range buff.clients {
		err := fn(client)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	buff.clients = buff.clients[:0]
	s.clientsPool.Put(buff)

	return nil
}

func (s *ClientStore) Channels(clientID UUID) ([]string, error) {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.channels: took=%s", end)
			}
		}()
	}

	clientsShard := s.clientsShard(clientID)

	_, channels, ok := clientsShard.Channels(clientID)
	if !ok {
		return nil, errors.WithStack(NewClientNotFoundError(clientID))
	}

	return channels, nil
}

func (s *ClientStore) CountChannels(clientID UUID) (int, error) {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.count_channels: took=%s", end)
			}
		}()
	}

	clientsShard := s.clientsShard(clientID)

	_, count, ok := clientsShard.CountChannels(clientID)
	if !ok {
		return 0, errors.WithStack(NewClientNotFoundError(clientID))
	}

	return count, nil
}

func (s *ClientStore) SetChannels(clientID UUID, channels ...string) error {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.set_channels: took=%s", end)
			}
		}()
	}

	if len(channels) == 0 {
		return nil
	}

	clientsShard := s.clientsShard(clientID)

	client, ok := clientsShard.SetChannels(clientID, channels...)
	if !ok {
		return errors.WithStack(NewClientNotFoundError(clientID))
	}

	for _, channel := range channels {
		channelsShard := s.channelsShard(channel)
		channelsShard.Link(client, channel)
	}

	return nil
}

func (s *ClientStore) UnsetChannels(clientID UUID, channels ...string) error {
	if s.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > s.options.DebugFuncTimeLimit {
				s.logger.Warnf("wspubsub.client_store.unset_channels: took=%s", end)
			}
		}()
	}

	clientsShard := s.clientsShard(clientID)

	_, ok := clientsShard.UnsetChannels(clientID, channels...)
	if !ok {
		return errors.WithStack(NewClientNotFoundError(clientID))
	}

	if len(channels) == 0 {
		for _, channelsShard := range s.channelsShardList {
			channelsShard.Unlink(clientID)
		}

		return nil
	}

	for _, channel := range channels {
		channelsShard := s.channelsShard(channel)
		channelsShard.Unlink(clientID, channel)
	}

	return nil
}

func (s *ClientStore) clientsShard(clientID UUID) *clientStoreClientsShard {
	index := xxhash.Sum64(clientID.Bytes()) % uint64(s.options.ClientShards.Count)

	return s.clientsShardList[index]
}

func (s *ClientStore) channelsShard(channel string) *clientStoreChannelsShard {
	index := xxhash.Sum64String(channel) % uint64(s.options.ChannelShards.Count)

	return s.channelsShardList[index]
}

func NewClientStore(options ClientStoreOptions, logger Logger) *ClientStore {
	clientList := &ClientStore{
		options:           options,
		logger:            logger,
		clientsShardList:  make([]*clientStoreClientsShard, options.ClientShards.Count),
		channelsShardList: make([]*clientStoreChannelsShard, options.ChannelShards.Count),
		clientsPool: sync.Pool{
			New: func() interface{} {
				return &clientsBuffer{clients: make([]WebsocketClient, 0, options.ClientShards.Size)}
			},
		},
	}

	for i := 0; i < options.ClientShards.Count; i++ {
		clientList.clientsShardList[i] = newClientStoreClientsShard(
			options.ClientShards.Size,
			options.ClientShards.BucketSize,
		)
	}

	for i := 0; i < options.ChannelShards.Count; i++ {
		clientList.channelsShardList[i] = newClientStoreChannelsShard(
			options.ChannelShards.Size,
			options.ChannelShards.BucketSize,
		)
	}

	return clientList
}
