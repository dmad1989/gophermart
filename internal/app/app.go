package app

import (
	"context"
)

type DB interface {
	Close() error
}

type App struct {
	db DB
}

func New(ctx context.Context, db DB) *App {
	return &App{db: db}
}
