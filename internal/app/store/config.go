package store

type Config struct {
	DatabaseURL string `toml:"database_url"`
	LogLevel    string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{}
}
