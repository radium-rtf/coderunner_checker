package domain

type Config struct {
	ServerConfig  `yaml:"server"`
	SandboxConfig `yaml:"sandbox"`
}

// possible env values
const (
	DevEnv  = "dev"
	ProdEnv = "prod"
	TestEnv = "test"
)

type ServerConfig struct {
	Env     string `yaml:"env" env-required:"true"`
	Address string `yaml:"address" env-required:"true"`
}

type Rule struct {
	Filename string `yaml:"filename" env-required:"true"`
	Image    string `yaml:"image" env-required:"true"`
	Launch   string `yaml:"launch" env-required:"true"`
}

type Rules map[string]Rule // Key is programming language, value is struct with specified launch rules

type SandboxConfig struct {
	User    string `yaml:"env" env-required:"true"`
	UUID    int    `yaml:"address" env-required:"true"`
	WorkDir string `yaml:"work_dir" env-required:"true"`
	Host    string `yaml:"host" env-required:"true"`
	Rules   Rules  `yaml:"rules" env-required:"true"`
}
