package conveyor

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dmad1989/gophermart/internal/client"
	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"github.com/golang/mock/gomock"
	"go.uber.org/goleak"
	"go.uber.org/zap"
)

type expectedOrdersCalc struct {
	err    error
	orders jsonobject.OrdersCalc
}

type expectedClientResponse struct {
	err      error
	accrual  jsonobject.AccrualResponse
	maxtimes int
}
type expectedUpdateResponse struct {
	err      error
	maxtimes int
}

func Test_conveyor_CalcProcess(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx := initContext()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mDB := NewMockDB(ctrl)
	mClient := NewMockClient(ctrl)

	tests := []struct {
		name          string
		expOrdersCalc expectedOrdersCalc
		expClient     expectedClientResponse
		expUpdate     expectedUpdateResponse
	}{
		{
			name:          "orderForCalc error",
			expOrdersCalc: expectedOrdersCalc{err: errors.New("")},
			expClient:     expectedClientResponse{err: nil, accrual: jsonobject.AccrualResponse{}, maxtimes: 0},
			expUpdate:     expectedUpdateResponse{err: nil, maxtimes: 0},
		},
		{
			name:          "ClientResponse fatal error",
			expOrdersCalc: expectedOrdersCalc{orders: make(jsonobject.OrdersCalc, 100)},
			expClient:     expectedClientResponse{err: client.ErrorAccrualFatal, accrual: jsonobject.AccrualResponse{}, maxtimes: 1},
			expUpdate:     expectedUpdateResponse{err: nil, maxtimes: 0},
		},
		{
			name:          "ClientResponse OverLoad error",
			expOrdersCalc: expectedOrdersCalc{orders: make(jsonobject.OrdersCalc, 100)},
			expClient:     expectedClientResponse{err: client.ErrorAccrualOverLoad, accrual: jsonobject.AccrualResponse{}, maxtimes: 1},
			expUpdate:     expectedUpdateResponse{err: nil, maxtimes: 0},
		},
		{
			name:          "ClientResponse UnknownOrder error",
			expOrdersCalc: expectedOrdersCalc{orders: make(jsonobject.OrdersCalc, 100)},
			expClient:     expectedClientResponse{err: client.ErrorAccrualUnknownOrder, accrual: jsonobject.AccrualResponse{}, maxtimes: 100},
			expUpdate:     expectedUpdateResponse{err: nil, maxtimes: 0},
		},
		{
			name:          "UpdateResponse times -1 batch",
			expOrdersCalc: expectedOrdersCalc{orders: make(jsonobject.OrdersCalc, 100)},
			expClient:     expectedClientResponse{accrual: jsonobject.AccrualResponse{Order: "11", Status: "NEW", Accrual: 100.0}, maxtimes: 100},
			expUpdate:     expectedUpdateResponse{err: nil, maxtimes: 100},
		},
		{
			name:          "UpdateResponse times -10 batch",
			expOrdersCalc: expectedOrdersCalc{orders: make(jsonobject.OrdersCalc, 1000)},
			expClient:     expectedClientResponse{accrual: jsonobject.AccrualResponse{Order: "11", Status: "NEW", Accrual: 100.0}, maxtimes: 1000},
			expUpdate:     expectedUpdateResponse{err: nil, maxtimes: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mDB.EXPECT().GetOrdersForCalc(gomock.Any()).Return(tt.expOrdersCalc.orders, tt.expOrdersCalc.err).MaxTimes(1)
			mClient.EXPECT().DoRequestAccrual(gomock.Any(), gomock.Any()).Return(tt.expClient.accrual, tt.expClient.err).MaxTimes(tt.expClient.maxtimes)
			mDB.EXPECT().UpdateOrders(gomock.Any(), gomock.Any()).Return(tt.expUpdate.err).MaxTimes(tt.expUpdate.maxtimes)
			c := conveyor{
				client: mClient,
				db:     mDB,
				logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}

			c.CalcProcess(ctx)
		})
	}
	time.Sleep(time.Second * 10)
	defer goleak.VerifyNone(t)
}

func loggerInit() (*zap.SugaredLogger, error) {
	zl, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("loggerInit: %w", err)
	}
	return zl.Sugar(), nil
}

func initContext() context.Context {
	log, err := loggerInit()
	if err != nil {
		log.Fatal(err)
	}
	return context.WithValue(context.Background(), config.LoggerCtxKey, log)
}
