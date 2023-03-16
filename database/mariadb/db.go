package mariadb

import (
	"context"
	"errors"

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

type UserRepo struct {
	DB *database.DB
}

func (u UserRepo) CreateUser(ctx context.Context, in database.NewUser) (*database.User, error) {
	return nil, errors.New("not implemented")
}

func (u UserRepo) ListUsers(ctx context.Context) ([]database.User, error) {
	return nil, errors.New("not implemented")
}
