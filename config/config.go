package config

type Config struct {
	DatabaseURL string
}

func LoadConfig(databaseUrl string) Config {
	return Config{
		DatabaseURL: databaseUrl,
	}
}
