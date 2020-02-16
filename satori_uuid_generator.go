package wspubsub

import uuid "github.com/satori/go.uuid"

var _ UUIDGenerator = (*SatoriUUIDGenerator)(nil)

// SatoriUUIDGenerator is an implementation of UUIDGenerator.
type SatoriUUIDGenerator struct{}

// GenerateV4 returns random generated UUID.
func (s SatoriUUIDGenerator) GenerateV4() UUID {
	return UUID(uuid.NewV4())
}
