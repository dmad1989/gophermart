package app

import (
	"context"
)

type DB interface {
	Close() error
	CreateOrder(ctx context.Context, orderNum uint64) error
	GetOrderAuthor(ctx context.Context, orderNum uint64) (int, error)
}

type App struct {
	db DB
}

func New(ctx context.Context, db DB) *App {
	return &App{db: db}
}

func (a App) CreateOrder(ctx context.Context, orderNum uint64) error {
	return nil
}
