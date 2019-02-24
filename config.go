package main

import (
	"errors"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

func loadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = toml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	err = checkRequired(&c)
	if err != nil {
		return nil, err
	}
	populateDefaults(&c)
	return &c, nil
}

func populateDefaults(c *Config) {
	if c.Redis.Address == "" {
		c.Redis.Address = "localhost:6379"
	}
}

func checkRequired(c *Config) error {
	if c.General.Token == "" {
		return errors.New("general.token is required")
	}
	if c.TdLib.Enabled && (c.TdLib.ApiID == 0 || c.TdLib.ApiHash == "") {
		return errors.New("tdlib.api_id and tdlib.api_hash are required if tdlib.enabled is set to true")
	}
	return nil
}

type Config struct {
	General GeneralConfig `toml:"general"`
	Redis RedisConfig `toml:"redis"`
	TdLib TdLibConfig `toml:"tdlib"`
}

type GeneralConfig struct {
	Token string `toml:"token"`
}

type TdLibConfig struct {
	Enabled bool `toml:"enabled"`
	ApiID int `toml:"api_id"`
	ApiHash string `toml:"api_hash"`
}

type RedisConfig struct {
	Address string `toml:"address"`
	Password string `toml:"password"`
	Database int `toml:"database"`
}