package conveyor

import (
	"context"
	"time"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"go.uber.org/zap"
)

type conveyor struct {
	logger *zap.SugaredLogger
	client Client
	db     DB
}
type Client interface {
	DoRequestAccrual(ctx context.Context, orderNum int) (jsonobject.AccrualResponse, error)
}
type DB interface {
}

func Start(ctx context.Context, client Client, db DB) {
	conv := conveyor{
		client: client,
		db:     db,
		logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}
	conv.logger.Infoln("worker start")
	tCh := time.NewTicker(time.Duration(time.Second * 30)).C
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
	c.logger.Infoln("worker in progress")
	//todo
	//c.db.findOredersToRequeest
	// c.client.DoRequestAccrual(ctx, orderNum)
	//  c.db.updateOrders
}
