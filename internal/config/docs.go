package config

type DocsConfig struct {
	Port         int    `yaml:"port" env-required:"true"`
	HttpENDPOINT string `yaml:"http_endpoint" env-required:"true"`
}
