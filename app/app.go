package app

import (
	"encoding/binary"

	"github.com/google/uuid"
)

type App struct {
	id uint64
}

func New() (*App, error) {
	appId, err := generateAppId()
	if err != nil {
		return nil, err
	}

	return &App{
		appId,
	}, nil
}

func generateAppId() (uint64, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return 0, err
	}

	// Create an app ID out of the first 8 bytes of the UUID.
	return binary.LittleEndian.Uint64(id[0:9]), nil
}
