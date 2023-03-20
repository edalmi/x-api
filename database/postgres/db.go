package postgres

import (
	"github.com/edalmi/x-api/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(dsn string) (*database.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &database.DB{
		DB: db,
	}, nil
}
