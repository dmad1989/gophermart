package wallet

import (
	"context"
	"net/http"

	"github.com/dmad1989/gophermart/internal/config"
	"go.uber.org/zap"
)

type App interface{}

type wallet struct {
	logger *zap.SugaredLogger
	app    App
}

func New(ctx context.Context, app App) *wallet {
	return &wallet{app: app, logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}
}

func (w wallet) PostOrdersHandler(res http.ResponseWriter, req *http.Request) {
	// req.Cookies()
	res.WriteHeader(http.StatusOK)
}

func (w wallet) GetOrdersHandler(res http.ResponseWriter, req *http.Request) {
	// req.Cookies()
	res.WriteHeader(http.StatusOK)
}

func (w wallet) BalanceHandler(res http.ResponseWriter, req *http.Request) {
	// req.Cookies()
	res.WriteHeader(http.StatusOK)
}

func (w wallet) WithdrawHandler(res http.ResponseWriter, req *http.Request) {
	// req.Cookies()
	res.WriteHeader(http.StatusOK)
}

func (w wallet) AllWithdrawalsHandler(res http.ResponseWriter, req *http.Request) {
	// req.Cookies()
	res.WriteHeader(http.StatusOK)
}
