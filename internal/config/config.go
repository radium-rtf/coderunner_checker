package config

import "github.com/radium-rtf/coderunner_checker/internal/domain"

// possible env values
const (
	LocalEnv = "local"
	DevEnv   = "dev"
	ProdEnv  = "prod"
)

type Config struct {
	Server  ServerConfig  `yaml:"server" env-required:"true"`
	Sandbox SandboxConfig `yaml:"sandbox" env-required:"true"`
	Env     string        `yaml:"env" env-required:"true"`
}

type ServerConfig struct {
	Port int `yaml:"port" env-required:"true"`
}

type SandboxConfig struct {
	User    string       `yaml:"user" env-default:"root"`
	UUID    int          `yaml:"uuid" env-default:"0"`
	WorkDir string       `yaml:"work_dir" env-required:"true"`
	Host    string       `yaml:"host" env-required:"true"`
	Rules   domain.Rules `yaml:"rules" env-required:"true"`
}
