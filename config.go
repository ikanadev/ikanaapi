package main

import "os"

type Config struct {
	Port             string
	DBConn           string
	MigrationsSource string
}

//nolint:gochecknoglobals
var _config = Config{
	Port:             os.Getenv("PORT"),
	DBConn:           os.Getenv("DATABASE"),
	MigrationsSource: os.Getenv("MIGRATIONS_SOURCE"),
}

func GetConfig() Config {
	return _config
}
