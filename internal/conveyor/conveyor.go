package conveyor

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dmad1989/gophermart/internal/client"
	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"go.uber.org/zap"
)

const batchSize = 100

var tikerTimeout time.Duration = 10
var sleepTime time.Duration = 0

type conveyor struct {
	logger *zap.SugaredLogger
	client Client
	db     DB
}
type Client interface {
	DoRequestAccrual(ctx context.Context, orderNum int) (jsonobject.AccrualResponse, error)
}
type DB interface {
	GetOrdersForCalc(ctx context.Context) (jsonobject.OrdersCalc, error)
	UpdateOrders(ctx context.Context, orders jsonobject.OrdersCalc) error
}

func Start(ctx context.Context, client Client, db DB) {
	conv := conveyor{
		client: client,
		db:     db,
		logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}
	conv.logger.Infoln("worker start")
	tCh := time.NewTicker(time.Duration(time.Second * tikerTimeout)).C
	go conv.doReapeat(ctx, tCh)
}

func (c conveyor) doReapeat(ctx context.Context, tCh <-chan time.Time) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Infoln("worker done")
			return
		case <-tCh:
			c.CalcProcess(ctx)
		}
	}
}

func (c conveyor) CalcProcess(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	c.logger.Infoln("worker in progress")
	orders, err := c.db.GetOrdersForCalc(ctx)
	if err != nil {
		c.logger.Infow("error in GetOrdersForCalc", zap.Error(err))
		cancel()
		return
	}
	c.logger.Infoln("step 1 done", zap.Int("len", len(orders)))

	var updOrders jsonobject.OrdersCalc
	bCh := make(chan jsonobject.OrdersCalc)
	defer close(bCh)
	go func(ctx context.Context, bCh chan jsonobject.OrdersCalc) {
		for b := range bCh {
			c.logger.Infoln("step 3 processing")
			err := c.db.UpdateOrders(ctx, b)
			if err != nil {
				c.logger.Infow("error in UpdateOrders", zap.Error(err))
				cancel()
				return
			}
		}
	}(ctx, bCh)

	start := 0

	for i, order := range orders {
		c.logger.Infoln(i)
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
		accrual, err := c.client.DoRequestAccrual(ctx, order.Number)
		if err != nil {
			c.logger.Infow("error in DoRequestAccrual", zap.Error(err))
			switch {
			case errors.Is(err, client.ErrorAccrualFatal):
				cancel()
				// close(bCh)
				return
			case errors.Is(err, client.ErrorAccrualOverLoad):
				// если сервис перегружен добавляем задержки между вызывами
				sleepTime = sleepTime + 1
				tikerTimeout = tikerTimeout + time.Duration(len(orders))*sleepTime
				cancel()
				// close(bCh)
				return
			case errors.Is(err, client.ErrorAccrualUnknownOrder):
				continue
			}

		}
		if !order.CalcStatus.Valid || accrual.Status != order.CalcStatus.String {
			c.logger.Infoln("step 2 processing", zap.Int("len updOrders", len(updOrders)))
			order.Accrual = accrual.Accrual
			order.CalcStatus = sql.NullString{String: accrual.Status, Valid: true}
			updOrders = append(updOrders, order)
			if len(orders) == i+1 || len(updOrders) == batchSize {
				select {
				case <-ctx.Done():
					return
				case bCh <- updOrders[start:]:
				}
			}
			start = len(updOrders)
		}
	}
	// close(bCh)
}
