package telnet

import (
	"time"
)

type Config struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func NewConfig(host, port string, timeout time.Duration) *Config {
	return &Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
}
