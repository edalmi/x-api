package database

import "github.com/jmoiron/sqlx"

type DB struct {
	*sqlx.DB
}

type User struct {
	ID string `db:"id"`
}

type NewUser struct {
	ID string `db:"id"`
}
