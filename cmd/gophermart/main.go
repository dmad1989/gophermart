package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmad1989/gophermart/internal/api"
	"github.com/dmad1989/gophermart/internal/app"
	"github.com/dmad1989/gophermart/internal/auth"
	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/db"
	"go.uber.org/zap"
)

func main() {
	log, err := loggerInit()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), config.LoggerCtxKey, log)
	defer log.Sync()
	conf := config.ParseConfig()
	db, err := db.New(ctx, conf.DbConnName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := app.New(ctx, db)
	auth := auth.New(ctx, db)
	api := api.New(ctx, app, conf.AccrualURL, auth)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	err = api.SeverStart(ctx, conf.ApiURL)
	if err != nil {
		panic(err)
	}
}

func loggerInit() (*zap.SugaredLogger, error) {
	zl, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("loggerInit: %w", err)
	}
	return zl.Sugar(), nil
}
