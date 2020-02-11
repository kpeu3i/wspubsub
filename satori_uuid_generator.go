package wspubsub

import uuid "github.com/satori/go.uuid"

var _ UUIDGenerator = (*SatoriUUIDGenerator)(nil)

type SatoriUUIDGenerator struct{}

func (s SatoriUUIDGenerator) GenerateV4() UUID {
	return UUID(uuid.NewV4())
}
