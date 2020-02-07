package wspubsub

import (
	"sync"
)

type clientStoreChannelsShardBucket map[UUID]WebsocketClient

type clientStoreChannelsShard struct {
	bucketSize int
	mu         sync.RWMutex
	clients    map[string]clientStoreChannelsShardBucket
}

func (s *clientStoreChannelsShard) Link(client WebsocketClient, channel string) {
	s.mu.Lock()
	if _, ok := s.clients[channel]; !ok {
		s.clients[channel] = make(clientStoreChannelsShardBucket, s.bucketSize)
	}
	s.clients[channel][client.ID()] = client
	s.mu.Unlock()
}

func (s *clientStoreChannelsShard) Unlink(clientID UUID, channels ...string) {
	if len(channels) == 0 {
		s.mu.Lock()
		for channel := range s.clients {
			delete(s.clients[channel], clientID)
		}
		s.mu.Unlock()

		return
	}

	s.mu.Lock()
	for _, channel := range channels {
		delete(s.clients[channel], clientID)
	}
	s.mu.Unlock()
}

func (s *clientStoreChannelsShard) Count(channels ...string) int {
	count := 0

	if len(channels) == 0 {
		s.mu.RLock()
		for channel := range s.clients {
			count += len(s.clients[channel])
		}
		s.mu.RUnlock()

		return count
	}

	s.mu.RLock()
	for _, channel := range channels {
		count += len(s.clients[channel])
	}
	s.mu.RUnlock()

	return count
}

func (s *clientStoreChannelsShard) Iterate(iterateFunc func(client WebsocketClient, channel string)) {
	s.mu.RLock()
	for channel, clients := range s.clients {
		for _, client := range clients {
			iterateFunc(client, channel)
		}
	}
	s.mu.RUnlock()
}

func newClientStoreChannelsShard(size int, bucketSize int) *clientStoreChannelsShard {
	shard := &clientStoreChannelsShard{
		bucketSize: bucketSize,
		clients:    make(map[string]clientStoreChannelsShardBucket, size),
	}

	return shard
}
