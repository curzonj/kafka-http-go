package main

import (
	"os"
)

const (
	DEFAULT_PORT = "3000"
)

type Config struct {
}

func (c *Config) HttpPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		value = DEFAULT_PORT
	}

	return value
}
