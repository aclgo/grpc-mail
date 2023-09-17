package config

import "time"

type Config struct {
	IntervalSend time.Duration
	Workers      int
	Redis        *Redis
	Gmail        *Gmail
	Ses          *Ses
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Ses struct {
}

type Gmail struct {
	Identity string
	Username string
	Password string
	Host     string
	Port     string
}

func NewConfig() *Config {
	return &Config{}
}
