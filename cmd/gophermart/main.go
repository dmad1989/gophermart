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
	"github.com/dmad1989/gophermart/internal/client"
	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/conveyor"
	"github.com/dmad1989/gophermart/internal/db"
	"github.com/dmad1989/gophermart/internal/gzipapi"
	"github.com/dmad1989/gophermart/internal/wallet"
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
	api := api.New(
		ctx,
		auth.New(ctx, db),
		gzipapi.New(ctx),
		wallet.New(
			ctx,
			app.New(ctx, db)))
	client := client.New(ctx, conf.AccrualURL)
	conveyor.Start(ctx, client, db)
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
