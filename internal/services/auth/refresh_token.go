package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RefreshToken struct {
	ID        uuid.UUID
	Hash      []byte
	ExpiresAt time.Time
}

func generateRefreshTokenBytes() ([]byte, error) {
	refreshUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "create uuid")
	}
	refreshBytes, err := refreshUUID.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "marshal uuid")
	}

	return refreshBytes, nil
}
