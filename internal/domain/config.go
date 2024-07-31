package domain

// possible env values
const (
	LocalEnv = "local"
	DevEnv   = "dev"
	ProdEnv  = "prod"
	TestEnv  = "test"
)

type Config struct {
	Server  ServerConfig  `yaml:"server" env-required:"true"`
	Sandbox SandboxConfig `yaml:"sandbox" env-required:"true"`
	Env     string        `yaml:"env" env-required:"true"`
}

type ServerConfig struct {
	Port int `yaml:"port" env-required:"true"`
}

type Rule struct {
	Filename string `yaml:"filename" env-required:"true"`
	Image    string `yaml:"image" env-required:"true"`
	Launch   string `yaml:"launch" env-required:"true"`
}

type Rules map[string]Rule // Key is programming language, value is struct with specified launch rules

type SandboxConfig struct {
	User    string `yaml:"user" env-required:"true"`
	UUID    int    `yaml:"uuid" env-required:"true"`
	WorkDir string `yaml:"work_dir" env-required:"true"`
	Host    string `yaml:"host" env-required:"true"`
	Rules   Rules  `yaml:"rules" env-required:"true"`
}
