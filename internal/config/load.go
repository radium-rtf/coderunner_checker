package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "config/config.yaml"

func Load() (*Config, error) {
	cfg := new(Config)

	path, ok := os.LookupEnv("CONFIG_PATH")
	if !ok || path == "" {
		path = configPath
	}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
