package rabbitmq

type Config struct {
	Url     string `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
	Retries int    `json:"retries" yaml:"retries" toml:"retries" mapstructure:"retries"`
}

func NewConfig() (conf *Config) {
	conf = new(Config)
	return
}
