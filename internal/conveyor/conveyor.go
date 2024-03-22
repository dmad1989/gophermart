package conveyor

import (
	"context"
	"database/sql"
	"time"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"go.uber.org/zap"
)

const batchSize = 100

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
	tCh := time.NewTicker(time.Duration(time.Second * 60)).C
	go conv.doReapeat(ctx, tCh)
}

func (c conveyor) doReapeat(ctx context.Context, tCh <-chan time.Time) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Infoln("worker done")
			return
		case <-tCh:
			c.work(ctx)
		}
	}
}

func (c conveyor) work(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	c.logger.Infoln("worker in progress")
	orders, err := c.db.GetOrdersForCalc(ctx)
	if err != nil {
		c.logger.Infow("error in GetOrdersForCalc", zap.Error(err))
		cancel()
		return
	}

	var updOrders jsonobject.OrdersCalc
	bCh := make(chan jsonobject.OrdersCalc)
	go func(ctx context.Context, bCh chan jsonobject.OrdersCalc) {
		for b := range bCh {
			err := c.db.UpdateOrders(ctx, b)
			if err != nil {
				c.logger.Infow("error in UpdateOrders", zap.Error(err))
				cancel()
				return
			}
		}
	}(ctx, bCh)

	for i, order := range orders {
		accrual, err := c.client.DoRequestAccrual(ctx, order.Number)
		if err != nil {
			c.logger.Infow("error in DoRequestAccrual", zap.Error(err))
			cancel()
			return
		}
		//TODO как то надо обнулять updOrders возможно сделать по другому
		if !order.CalcStatus.Valid || accrual.Status != order.CalcStatus.String {
			order.Accrual = accrual.Accrual
			order.CalcStatus = sql.NullString{String: accrual.Status, Valid: true}
			updOrders = append(updOrders, order)
			if len(orders) == i+1 || len(updOrders) == batchSize {
				select {
				case <-ctx.Done():
					return
				case bCh <- updOrders:
				}
			}
		}
	}
	close(bCh)
}
