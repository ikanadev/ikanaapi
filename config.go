package main

import "os"

type Config struct {
	Port          string
	DBConn        string
	MigrationsURL string
}

//nolint:gochecknoglobals
var _config = Config{
	Port:          os.Getenv("PORT"),
	DBConn:        os.Getenv("DATABASE"),
	MigrationsURL: os.Getenv("MIGRATIONS_URL"),
}

func GetConfig() Config {
	return _config
}

