package centrifugo

type Config struct {
	Url    string `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
	ApiKey string `json:"api_key" yaml:"api_key" toml:"api_key" mapstructure:"api_key"`
}

func NewConfig() (conf *Config) {
	conf = new(Config)
	return
}
