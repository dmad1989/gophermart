package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/dmad1989/gophermart/internal/jsonobject"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

const timeout = time.Duration(time.Second * 10)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

//go:embed sql/insertUser.sql
var sqlInsertUser string

//go:embed sql/getUserPassword.sql
var sqlUserPassword string

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

func (db *DB) CreateUser(ctx context.Context, user jsonobject.User) error {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := db.conn.NamedExecContext(tctx, sqlInsertUser, user)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (db *DB) GetUserPassword(ctx context.Context, login string) (string, error) {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := db.conn.NamedExecContext(tctx, sqlInsertUser, user)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}
