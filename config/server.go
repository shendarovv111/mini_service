package config

import (
	"os"
	"strconv"
)

type ServerConfig struct {
	Port int
}

func NewServerConfig() *ServerConfig {
	port := 8080

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil && p > 0 {
			port = p
		}
	}

	return &ServerConfig{Port: port}
}
