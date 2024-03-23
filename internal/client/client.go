package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"go.uber.org/zap"
)

const URLPattern = "%s/api/orders/%d"

var (
	ErrorAccrualFatal        = errors.New("accrual service is unavailable")
	ErrorAccrualUnknownOrder = errors.New("order not registed in accrual service")
	ErrorAccrualOverLoad     = errors.New("accrual service is overload")
)

type client struct {
	logger     *zap.SugaredLogger
	accrualURL string
}

func New(ctx context.Context, accrualURL string) *client {
	return &client{logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger), accrualURL: accrualURL}
}

func (c *client) DoRequestAccrual(ctx context.Context, orderNum int) (jsonobject.AccrualResponse, error) {
	accRes := jsonobject.AccrualResponse{}
	defer c.logger.Sync()
	u := fmt.Sprintf(URLPattern, c.accrualURL, orderNum)
	c.logger.Infow("accrual request", zap.String("request url", u))
	res, err := http.Get(u)
	if err != nil {
		c.logger.Infow("in accrual request", zap.String("error", err.Error()))
		return accRes, fmt.Errorf("accrual request: %w", err)
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return accRes, fmt.Errorf("accrual request: reading response body: %w", err)
		}
		if err := accRes.UnmarshalJSON(body); err != nil {
			return accRes, fmt.Errorf("accrual request: decoding response: %w", err)
		}
		return accRes, nil
	case http.StatusNoContent:
		return accRes, ErrorAccrualUnknownOrder
	case http.StatusTooManyRequests:
		return accRes, ErrorAccrualOverLoad
	default:
		return accRes, ErrorAccrualFatal
	}
}
