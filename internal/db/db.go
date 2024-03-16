package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

const timeout = time.Duration(time.Second * 10)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

//go:embed sql/insertUser.sql
var sqlInsertUser string

//go:embed sql/getUserByLogin.sql
var sqlUserByLogin string

type DB struct {
	conn   *sqlx.DB
	logger *zap.SugaredLogger
}

func New(ctx context.Context, connName string) (*DB, error) {
	if connName == "" {
		return nil, errors.New("new db: empty conn Name")
	}
	conn, err := sqlx.Connect("pgx", connName)
	if err != nil {
		return nil, fmt.Errorf("conncet to DB: %w", err)
	}
	res := DB{conn: conn, logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(conn.DB, "sql/migrations"); err != nil {
		return nil, fmt.Errorf("goose: create table: %w", err)
	}
	res.logger.Infow("db started!")
	return &res, nil
}

func (db *DB) Close() error {
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("close db conn: %w", err)
	}
	return nil
}

func (db *DB) CreateUser(ctx context.Context, user jsonobject.User) (int, error) {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	db.logger.Infow("Creating user",
		zap.String("login", user.Login),
		zap.String("password", user.Password),
		zap.ByteString("hashed", user.HashPassword))

	res, err := db.conn.NamedExecContext(tctx, sqlInsertUser, user)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fetching LastInsertId: %w", err)
	}
	return int(id), nil
}

func (db *DB) GetUserByLogin(ctx context.Context, login string) (jsonobject.User, error) {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	user := jsonobject.User{}
	err := db.conn.GetContext(tctx, &user, sqlUserByLogin, login)
	if err != nil {
		return jsonobject.User{}, fmt.Errorf("GetUserByLogin: %w", err)
	}
	return user, nil
}

func (db *DB) CheckUserExists(ctx context.Context, login string) (bool, error) {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	user := jsonobject.User{}
	err := db.conn.GetContext(tctx, &user, sqlUserByLogin, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("CheckUserExists: %w", err)
	}
	return true, nil
}
