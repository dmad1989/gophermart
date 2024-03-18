package wallet

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmad1989/gophermart/internal/config"
	"go.uber.org/zap"
)

var (
	ErrorRequestContentType = errors.New("wrong content-type")
	ErrorRequestEmptyBody   = errors.New("empty body not expected")
)

type App interface {
	CreateOrder(ctx context.Context, orderNum uint64) error
}

type wallet struct {
	logger *zap.SugaredLogger
	app    App
}

func New(ctx context.Context, app App) *wallet {
	return &wallet{app: app, logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}
}

func (w wallet) PostOrdersHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
		errorResponse(res, http.StatusBadRequest, ErrorRequestContentType)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("reading request body: %w", err))
		return
	}

	if len(body) <= 0 {
		errorResponse(res, http.StatusBadRequest, ErrorRequestEmptyBody)
		return
	}

	err = w.app.CreateOrder(req.Context(), binary.BigEndian.Uint64(body))

	// if err == "номер заказа уже был загружен этим пользователем;" {
	// errorResponse(res, http.StatusOK, ErrorRequestContentType)
	// res.WriteHeader(http.StatusOK)
	// res.Write([]byte("номер заказа уже был загружен этим пользователем"))
	// return
	// }
	// if err == "неверный формат номера заказа;" {
	// res.WriteHeader(http.StatusUnprocessableEntity)
	// res.Write([]byte("неверный формат номера заказа"))
	// return
	// }
	// if err == "номер заказа уже был загружен другим пользователем;" {
	// res.WriteHeader(http.StatusConflict)
	// res.Write([]byte("неверный формат номера заказа"))
	// return
	// }
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("create order: %w", err))
		return
	}

	res.WriteHeader(http.StatusAccepted)
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

func errorResponse(res http.ResponseWriter, status int, err error) {
	res.WriteHeader(status)
	res.Write([]byte(err.Error()))
}
