package wspubsub

import (
	"sync"

	"github.com/pkg/errors"
)

type clientStoreShardValue struct {
	client   WebsocketClient
	channels map[string]struct{}
}

type clientStoreClientsShard struct {
	bucketSize int
	mu         sync.RWMutex
	values     map[UUID]clientStoreShardValue
}

func (s *clientStoreClientsShard) Get(clientID UUID) (WebsocketClient, error) {
	s.mu.RLock()
	value, ok := s.values[clientID]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.WithStack(NewClientNotFoundError(clientID))
	}

	return value.client, nil
}

func (s *clientStoreClientsShard) Set(client WebsocketClient) {
	s.mu.Lock()
	s.values[client.ID()] = clientStoreShardValue{client: client, channels: make(map[string]struct{}, s.bucketSize)}
	s.mu.Unlock()
}

func (s *clientStoreClientsShard) Unset(clientID UUID) {
	s.mu.Lock()
	delete(s.values, clientID)
	s.mu.Unlock()
}

func (s *clientStoreClientsShard) Count() int {
	s.mu.RLock()
	count := len(s.values)
	s.mu.RUnlock()

	return count
}

func (s *clientStoreClientsShard) Channels(clientID UUID) (WebsocketClient, []string, bool) {
	var (
		client   WebsocketClient
		channels []string
	)

	s.mu.RLock()
	value, ok := s.values[clientID]
	if ok {
		client = value.client
		channels = make([]string, 0, len(value.channels))
		for channel := range value.channels {
			channels = append(channels, channel)
		}
	}
	s.mu.RUnlock()

	return client, channels, ok
}

func (s *clientStoreClientsShard) CountChannels(clientID UUID) (WebsocketClient, int, bool) {
	var (
		client WebsocketClient
		count  int
	)

	s.mu.RLock()
	value, ok := s.values[clientID]
	if ok {
		client = value.client
		count = len(value.channels)
	}
	s.mu.RUnlock()

	return client, count, ok
}

func (s *clientStoreClientsShard) SetChannels(clientID UUID, channels ...string) (WebsocketClient, bool) {
	var client WebsocketClient

	s.mu.Lock()
	value, ok := s.values[clientID]
	if ok {
		client = value.client
		for _, channel := range channels {
			value.channels[channel] = struct{}{}
		}
	}
	s.mu.Unlock()

	return client, ok
}

func (s *clientStoreClientsShard) UnsetChannels(clientID UUID, channels ...string) (WebsocketClient, bool) {
	var client WebsocketClient

	s.mu.Lock()
	value, ok := s.values[clientID]
	if ok {
		client = value.client
		if len(channels) == 0 {
			for channel := range value.channels {
				delete(value.channels, channel)
			}
		} else {
			for _, channel := range channels {
				delete(value.channels, channel)
			}
		}
	}
	s.mu.Unlock()

	return client, ok
}

func (s *clientStoreClientsShard) Iterate(iterateFunc func(client WebsocketClient)) {
	s.mu.RLock()
	for _, value := range s.values {
		iterateFunc(value.client)
	}
	s.mu.RUnlock()
}

func newClientStoreClientsShard(size int, bucketSize int) *clientStoreClientsShard {
	return &clientStoreClientsShard{bucketSize: bucketSize, values: make(map[UUID]clientStoreShardValue, size)}
}
