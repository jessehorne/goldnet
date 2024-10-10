package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type ServerConfig struct {
	WorldSeed     int64
	ServerAddress string
}

func NewServerConfig() (*ServerConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	seed := os.Getenv("WORLD_SEED")
	if seed == "" {
		return nil, errors.New("missing WORLD_SEED")
	}
	seedNum, err := strconv.Atoi(seed)
	if err != nil {
		return nil, errors.New("couldn't parse WORLD_SEED:" + err.Error())
	}

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		return nil, errors.New("missing SERVER_ADDRESS")
	}

	return &ServerConfig{
		WorldSeed:     int64(seedNum),
		ServerAddress: addr,
	}, nil
}
