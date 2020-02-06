package wspubsub

import uuid "github.com/satori/go.uuid"

type SatoriUUIDGenerator struct{}

func (s SatoriUUIDGenerator) GenerateV4() UUID {
	return UUID(uuid.NewV4())
}
