package wspubsub

import "time"

// ClientStoreOptions represents configuration of the storage.
type ClientStoreOptions struct {
	ClientShards struct {
		// Total number of shards
		Count int

		// Size of shard
		Size int

		// Size of a bucket in shard
		BucketSize int
	}
	ChannelShards struct {
		// Total number of shards
		Count int

		// Size of shard
		Size int

		// Size of a bucket in shard
		BucketSize int
	}

	// Enable/disable debug mode.
	IsDebug bool

	// Function execution time limit in debug mode.
	// Exceeding this time limit will cause a new warn log message.
	DebugFuncTimeLimit time.Duration
}

// NewClientStoreOptions initializes a new ClientStoreOptions.
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
