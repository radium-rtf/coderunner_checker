package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/radium-rtf/coderunner_checker/internal/domain"
)

const configPath = "config/config.yaml"

func Load() (*domain.Config, error) {
	cfg := new(domain.Config)

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
