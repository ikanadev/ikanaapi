package config

import "os"

type Config struct {
	Port             string
	DBConn           string
	MigrationsSource string
	OpenAIKey        string
}

var _config = Config{
	Port:             os.Getenv("PORT"),
	DBConn:           os.Getenv("DATABASE"),
	MigrationsSource: os.Getenv("MIGRATIONS_SOURCE"),
	OpenAIKey:        os.Getenv("OPENAI_API_KEY"),
}

func GetConfig() Config {
	return _config
}
