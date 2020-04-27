package main

import (
	"github.com/caarlos0/env"
)

func loadConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Token        string `env:"TG_TOKEN"`
	TdlibEnabled bool   `env:"TDLIB_ENABLED"`
	TdlibApiID   int    `env:"TDLIB_API_ID"`
	TdlibApiHash string `env:"TDLIB_API_HASH"`
	RedisAddr    string `env:"DB_ADDR"`
}
