package main

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	KafkaMaxBufferTime    uint32
	KafkaMaxBufferedBytes uint32
	KafkaBrokers          []string
	HttpPort              string
}

func (c *Config) Configure() {
	c.HttpPort = envVarOrDefaultString("PORT", "3000")
	c.KafkaMaxBufferTime = uint32(envVarOrDefaultInt("MAX_BUFFER_TIME", 2))
	c.KafkaMaxBufferedBytes = uint32(envVarOrDefaultInt("MAX_BUFFER_BYTES", 16000))
	c.KafkaBrokers = strings.Split(envVarOrDefaultString("BROKERS", "localhost:9092"), ",")
}

func envVarOrDefaultString(name string, value string) string {
	v := os.Getenv(name)
	if v == "" {
		v = value
	}

	return v
}

func envVarOrDefaultInt(name string, value int64) int64 {
	var v int64

	vs := os.Getenv(name)
	if vs == "" {
		v = value
	} else {
		var err error
		v, err = strconv.ParseInt(vs, 10, 64)
		if err != nil {
			panic(err)
		}
	}

	return v
}
