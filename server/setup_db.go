package server

import (
	"errors"

	"github.com/edalmi/x-api/config"
	"github.com/edalmi/x-api/database"
	"github.com/edalmi/x-api/database/mariadb"
	"github.com/edalmi/x-api/database/mysql"
	"github.com/edalmi/x-api/database/postgres"
	"github.com/edalmi/x-api/database/sqlite"
)

func setupDB(cfg *config.DB) (*database.DB, error) {
	if cfg.Postgres != nil {
		return postgres.New(cfg.Postgres.GetDSN())
	}

	if cfg.SQLite != nil {
		return sqlite.New(cfg.SQLite.GetDSN())
	}

	if cfg.MySQL != nil {
		return mysql.New(cfg.MySQL.GetDSN())
	}

	if cfg.MariaDB != nil {
		return mariadb.New(cfg.MariaDB.GetDSN())
	}

	return nil, errors.New("no database configuration found")
}
