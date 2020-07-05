package config

import (
	"fmt"
	"os"
)

const (
	DATABASE_USERNAME = "DATABASE_USERNAME"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_NAME     = "DATABASE_NAME"
)

type DBConfig struct {
	username string
	password string
	database string
}

func NewDbConfig() *DBConfig {
	return &DBConfig{
		username: os.Getenv(DATABASE_USERNAME),
		password: os.Getenv(DATABASE_PASSWORD),
		database: os.Getenv(DATABASE_NAME),
	}
}

func (dbc *DBConfig) String() string {
	return fmt.Sprintf("mongodb://%s:%s@127.0.0.1:27017/%s",
		dbc.username,
		dbc.password,
		dbc.database)
}

func (dbc *DBConfig) DatabaseName() string {
	return dbc.database
}
