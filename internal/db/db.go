package db

import (
	"context"
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

var (
	//go:embed sql/migrations/*.sql
	embedMigrations embed.FS

	//go:embed sql/insertUser.sql
	sqlInsertUser string

	//go:embed sql/getUserByLogin.sql
	sqlUserByLogin string

	//go:embed sql/insertOrder.sql
	sqlInsertOrder string

	//go:embed sql/getOrderAuthor.sql
	sqlCheckOrderAuthor string
)

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
	id := []int{}
	//Сделано через select так как exec возваращает sql.Result, у него есть lastInserted - но это не поддерживается в Postgres
	err := db.conn.SelectContext(tctx, &id, sqlInsertUser, user.Login, user.HashPassword)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	return id[0], nil
}

func (db *DB) GetUserByLogin(ctx context.Context, login string) (jsonobject.User, error) {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	user := jsonobject.User{}
	err := db.conn.GetContext(tctx, &user, sqlUserByLogin, login)
	if err != nil {
		return jsonobject.User{}, fmt.Errorf("db: GetUserByLogin: %w", err)
	}
	return user, nil
}

func (db *DB) CreateOrder(ctx context.Context, orderNum uint64) error {
	userID := ctx.Value(config.UserCtxKey)
	if userID == "" {
		return errors.New("db: createOrder: no user in context")
	}
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := db.conn.ExecContext(tctx, sqlInsertOrder, orderNum, userID)
	if err != nil {
		return fmt.Errorf("db: createOrder: %w", err)
	}
	return nil
}

func (db *DB) GetOrderAuthor(ctx context.Context, orderNum uint64) (int, error) {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	authorId := 0
	err := db.conn.SelectContext(tctx, &authorId, sqlCheckOrderAuthor, orderNum)
	if err != nil {
		return 0, fmt.Errorf("db: getOrderAuthor: %w", err)
	}
	return authorId, nil
}
