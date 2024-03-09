package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	conn *sqlx.DB
}

func New(ctx context.Context, connName string) (*DB, error) {

	if connName == "" {
		return nil, errors.New("new db: empty conn Name")
	}
	conn, err := sqlx.Connect("pgx", connName)
	if err != nil {
		return nil, fmt.Errorf("conncet to DB: %w", err)
	}

	res := DB{conn: conn}

	return &res, nil
}

func (db *DB) Close() error {
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("close db conn: %w", err)
	}
	return nil
}
