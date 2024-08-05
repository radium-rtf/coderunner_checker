package config

type Config struct {
	Env     string        `yaml:"env" env-required:"true"`
	Server  ServerConfig  `yaml:"server"`
	Sandbox SandboxConfig `yaml:"sandbox"`
}

// possible env values
const (
	DevEnv  = "dev"
	ProdEnv = "prod"
	TestEnv = "test"
)

type ServerConfig struct {
	Address string `yaml:"address" env-required:"true"`
}

type Rule struct {
	Filename string `yaml:"filename" env-required:"true"`
	Image    string `yaml:"image" env-required:"true"`
	Launch   string `yaml:"launch" env-required:"true"`
}

type Std struct {
	Versions map[string]Rule `yaml:"versions" env-required:"true"`
}

type Rules struct {
	Languages map[string]Std `yaml:"languages" env-required:"true"`
}

type SandboxConfig struct {
	User    string `yaml:"user" env-required:"true"`
	UUID    int    `yaml:"uuid" env-required:"true"`
	WorkDir string `yaml:"work_dir" env-required:"true"`
	Host    string `yaml:"host" env-required:"true"`
	Rules   Rules  `yaml:"rules" env-required:"true"`
}
