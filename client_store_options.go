package wspubsub

import "time"

type ClientStoreOptions struct {
	ClientShards struct {
		Count      int
		Size       int
		BucketSize int
	}
	ChannelShards struct {
		Count      int
		Size       int
		BucketSize int
	}
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// nolint: gomnd
func NewClientStoreOptions() ClientStoreOptions {
	options := ClientStoreOptions{
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}

	options.ClientShards.Count = 128
	options.ClientShards.Size = 10000
	options.ChannelShards.BucketSize = 100

	options.ChannelShards.Count = 16
	options.ChannelShards.Size = 100
	options.ChannelShards.BucketSize = 10000

	return options
}
