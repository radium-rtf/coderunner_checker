package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "config/config.yaml"

func Load() (*Config, error) {
	cfg := new(Config)

	path, ok := os.LookupEnv("CONFIG_PATH")
	if !ok || path == "" {
		return nil, fmt.Errorf("no config path in env")
	}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadFromPath(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadDocs() (*DocsConfig, error) {
	cfg := new(DocsConfig)

	path, ok := os.LookupEnv("DOCS_CONFIG_PATH")
	if !ok || path == "" {
		return nil, fmt.Errorf("no docs config path in env")
	}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
