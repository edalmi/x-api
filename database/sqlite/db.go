package sqlite

import (
	"github.com/edalmi/x-api/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func New(dsn string) (*database.DB, error) {
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return &database.DB{
		DB: db,
	}, nil
}
