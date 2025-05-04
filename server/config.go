package server

import "fmt"

type Config struct {
	Host     string
	Port     int
	LogLevel string
}

func NewConfig(host string, port int) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		LogLevel: "info",
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) SetLogLevel(level string) {
	c.LogLevel = level
}
