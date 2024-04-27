package models

type Config struct {
	Token  string `json:"token,omitempty" yaml:"token"`
	ApiKey string `json:"apiKey,omitempty" yaml:"apiKey"`
	Output string `json:"output,omitempty" yaml:"output"`
}

func NewConfig() *Config {
	return &Config{
		Token:  "",
		ApiKey: "",
		Output: "json",
	}
}
