package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmad1989/gophermart/internal/api"
	"github.com/dmad1989/gophermart/internal/app"
	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/db"
)

func main() {
	ctx := context.Background()
	conf := config.ParseConfig()
	db, err := db.New(ctx, conf.DbConnName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := app.New(ctx, db)

	api := api.New(ctx, app, conf.AccrualURL)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	err = api.SeverStart(ctx, conf.ApiURL)
	if err != nil {
		panic(err)
	}
}
