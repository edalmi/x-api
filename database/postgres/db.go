package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *DB {
	return &DB{
		DB: db,
	}
}

type DB struct {
	*sqlx.DB
}

func (d *DB) Open() error {
	return errors.New("not implemented")
}

func (d *DB) Close() error {
	return errors.New("not implemented")
}
