package config

import (
	"os"
	"strconv"
)

type Config struct {
	ListenHost string
	ListenPort int
	StoreType  string
}

func New() *Config {
	cfg := &Config{
		ListenHost: "0.0.0.0",
		ListenPort: 8080,
		StoreType:  "memory",
	}

	// Read env vars
	if host := os.Getenv("LISTEN_HOST"); host != "" {
		cfg.ListenHost = host
	}
	if port := os.Getenv("LISTEN_PORT"); port != "" {
		if portNum, err := strconv.Atoi(port); err == nil {
			cfg.ListenPort = portNum
		}
	}
	if storeType := os.Getenv("STORE_TYPE"); storeType != "" {
		cfg.StoreType = storeType
	}

	// Check storetype
	if cfg.StoreType != "memory" { // Add more store types as needed
		panic("Unsupported store type: " + cfg.StoreType)
	}

	return cfg
}
