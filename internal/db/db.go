package db

import (
	"context"
	"embed"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

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

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(conn.DB, "sql/migrations"); err != nil {
		return nil, fmt.Errorf("goose: create table: %w", err)
	}

	return &res, nil
}

func (db *DB) Close() error {
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("close db conn: %w", err)
	}
	return nil
}
