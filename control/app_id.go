package control

import (
	"encoding/binary"

	"github.com/google/uuid"
)

// NewAppId returns a random app ID as an uint64 number.
func NewAppId() uint64 {
	id, err := uuid.NewRandom()
	if err != nil {
		// Should never happen.
		panic(err)
	}

	return binary.LittleEndian.Uint64(id[0:9])
}
