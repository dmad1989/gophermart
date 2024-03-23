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
	ErrorFromatNumber    = errors.New("неверный формат номера заказа")
	ErrorOrderAuthor     = errors.New("заказ с таким же номером уже загружен другим пользователем")
	ErrorOrderUnique     = errors.New("заказ с таким же номером уже загружен")
	ErrorNotEnoughPoints = errors.New("недостаточно баллов для списания")
)

type DB interface {
	Close() error
	CreateOrder(ctx context.Context, orderNum int) error
	GetOrderAuthor(ctx context.Context, orderNum int) (int, error)
	GetOrdersByUser(ctx context.Context) (jsonobject.Orders, error)
	GetUserBalance(ctx context.Context) (jsonobject.Balance, error)
	CreateWithdraw(ctx context.Context, w jsonobject.Withdraw) error
	GetWithdrawlsByUser(ctx context.Context) (jsonobject.Withdrawls, error)
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
	authorID, err := a.db.GetOrderAuthor(ctx, orderNum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("app (createOrder): %w ", err)
	}

	if authorID != 0 {
		userID := ctx.Value(config.UserCtxKey)
		if userID == "" {
			return errors.New("app (createOrder): no user in context")
		}
		if authorID != userID {
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

func (a App) GetUserBalance(ctx context.Context) (jsonobject.Balance, error) {
	b, err := a.db.GetUserBalance(ctx)
	if err != nil {
		return jsonobject.Balance{}, fmt.Errorf("app (GetUserBalance):  %w", err)
	}

	b.Withdrawn = getValidValue(b.WithdrawnDB)
	b.AccrualCurrent = getValidValue(b.AccrualDB)

	if b.AccrualCurrent == 0 && b.Withdrawn > 0 {
		return jsonobject.Balance{}, errors.New("app (GetUserBalance): нет начислений, но есть списания")
	}
	//AccrualDB хранит в себе все когда либо начисленные баллы, чтобы узнать актуальный баланс вычитаем
	b.AccrualCurrent = b.AccrualCurrent - b.Withdrawn

	if b.AccrualCurrent < 0 {
		return jsonobject.Balance{}, errors.New("app (GetUserBalance): минусовой баланс счета")
	}

	return b, nil
}

func getValidValue(num sql.NullFloat64) float32 {
	var res float32
	if num.Valid {
		res = float32(num.Float64)
	}
	return res
}

func (a App) CreateWithdraw(ctx context.Context, w jsonobject.Withdraw) error {
	if !validateNumber(w.OrderNum) {
		return fmt.Errorf("app (CreateWithdraw): %w ", ErrorFromatNumber)
	}

	balance, err := a.GetUserBalance(ctx)
	if err != nil {
		return fmt.Errorf("app (CreateWithdraw): %w", err)
	}

	if balance.AccrualCurrent-w.Sum < 0 {
		return ErrorNotEnoughPoints
	}

	err = a.db.CreateWithdraw(ctx, w)
	if err != nil {
		return fmt.Errorf("app (CreateWithdraw): %w ", err)
	}
	return nil
}

func (a App) GetWithdrawlsByUser(ctx context.Context) (jsonobject.Withdrawls, error) {
	return a.db.GetWithdrawlsByUser(ctx)
}
