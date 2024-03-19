package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
)

var (
	ErrorFromatNumber = errors.New("неверный формат номера заказа")
	ErrorOrderAuthor  = errors.New("заказ с таким же номером уже загружен другим пользователем")
	ErrorOrderUnique  = errors.New("заказ с таким же номером уже загружен")
)

type DB interface {
	Close() error
	CreateOrder(ctx context.Context, orderNum int) error
	GetOrderAuthor(ctx context.Context, orderNum int) (int, error)
	GetOrdersByUser(ctx context.Context) (jsonobject.Orders, error)
}

type App struct {
	db DB
}

func New(ctx context.Context, db DB) *App {
	return &App{db: db}
}

func (a App) CreateOrder(ctx context.Context, orderNum int) error {
	if !validateNumber(orderNum) {
		return fmt.Errorf("app (createOrder): %w ", ErrorFromatNumber)
	}
	authorId, err := a.db.GetOrderAuthor(ctx, orderNum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("app (createOrder): %w ", err)
	}

	if authorId != 0 {
		userID := ctx.Value(config.UserCtxKey)
		if userID == "" {
			return errors.New("app (createOrder): no user in context")
		}
		if authorId != userID {
			return fmt.Errorf("app (createOrder): %w ", ErrorOrderAuthor)
		}
		return fmt.Errorf("app (createOrder): %w ", ErrorOrderUnique)
	}

	err = a.db.CreateOrder(ctx, orderNum)
	if err != nil {
		return fmt.Errorf("App (create order): %w", err)
	}

	return nil
}

func validateNumber(num int) bool {
	return (num%10+checksum(num/10))%10 == 0
}

func checksum(num int) int {
	var luhn int
	for i := 0; num > 0; i++ {
		cur := num % 10
		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}
		luhn += cur
		num = num / 10
	}
	return luhn % 10
}

func (a App) GetOrdersByUser(ctx context.Context) (jsonobject.Orders, error) {
	return a.db.GetOrdersByUser(ctx)
}
