package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/dmad1989/gophermart/internal/app"
	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"go.uber.org/zap"
)

var (
	ErrorRequestContentType   = errors.New("wrong content-type")
	ErrorRequestEmptyBody     = errors.New("empty body not expected")
	ErrorRequestContextNoUser = errors.New("no user in context")
)

type App interface {
	CreateOrder(ctx context.Context, orderNum int) error
	GetOrdersByUser(ctx context.Context) (jsonobject.Orders, error)
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

	orderNum, err := strconv.Atoi(string(body))
	if len(body) <= 0 {
		errorResponse(res, http.StatusBadRequest, fmt.Errorf("converting body to int: %w", err))
		return
	}

	err = w.app.CreateOrder(req.Context(), orderNum)

	if errors.Is(err, app.ErrorOrderUnique) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(err.Error()))
		return
	}
	if errors.Is(err, app.ErrorFromatNumber) {
		errorResponse(res, http.StatusUnprocessableEntity, fmt.Errorf("post order: %w", err))
		return
	}
	if errors.Is(err, app.ErrorOrderAuthor) {
		errorResponse(res, http.StatusConflict, fmt.Errorf("post order: %w", err))
		return
	}
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("create order: %w", err))
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

func (w wallet) GetOrdersHandler(res http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(config.UserCtxKey)
	if userID == nil || userID == 0 {
		errorResponse(res, http.StatusUnauthorized, ErrorRequestContextNoUser)
		return
	}
	orders, err := w.app.GetOrdersByUser(req.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("getOrders: %w", err))
		return
	}
	if len(orders) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	ordersJson, err := orders.MarshalJSON()
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("getOrders: encoding response: %w", err))
		return
	}
	res.Write(ordersJson)
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
