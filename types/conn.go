package types

import "github.com/fossmedaddy/dbdaddy/constants"

type ConnConfig struct {
	User     string            `json:"user"`
	Driver   string            `json:"driver"`
	Password string            `json:"password"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Database string            `json:"dbname"`
	Params   map[string]string `json:"params"`
}

func NewDefaultPgConnConfig() ConnConfig {
	return ConnConfig{
		User:     "postgres",
		Password: "postgres",
		Host:     "127.0.0.1",
		Port:     "5432",
		Database: "postgres",
		Driver:   constants.DbDriverPostgres,
		Params:   map[string]string{},
	}
}
