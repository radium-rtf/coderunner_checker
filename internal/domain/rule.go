package domain

type Rule struct {
	Filename  string `yaml:"filename" env-required:"true"`
	Image     string `yaml:"image" env-required:"true"`
	Launch    string `yaml:"launch" env-required:"true"`
}

type Rules map[string]Rule // Key is programming language, value is struct with specified launch rules
