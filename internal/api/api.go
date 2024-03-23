package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Wallet interface {
	PostOrdersHandler(res http.ResponseWriter, req *http.Request)
	GetOrdersHandler(res http.ResponseWriter, req *http.Request)
	BalanceHandler(res http.ResponseWriter, req *http.Request)
	WithdrawHandler(res http.ResponseWriter, req *http.Request)
	GetWithdrawalsHandler(res http.ResponseWriter, req *http.Request)
}
type Auth interface {
	LoginHandler(http.ResponseWriter, *http.Request)
	RegisterHandler(http.ResponseWriter, *http.Request)
	CheckMiddleware(h http.Handler) http.Handler
}
type gzipAPI interface {
	Middleware(h http.Handler) http.Handler
}

type api struct {
	logger *zap.SugaredLogger
	router *chi.Mux
	auth   Auth
	gzip   gzipAPI
	wallet Wallet
}

func New(ctx context.Context, auth Auth, gzip gzipAPI, wallet Wallet) *api {
	api := &api{
		router: chi.NewRouter(),
		auth:   auth,
		gzip:   gzip,
		wallet: wallet,
		logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger),
	}
	api.initRouter()
	return api
}

func (a api) initRouter() {
	a.router.Use(middleware.Logger, middleware.Recoverer, a.gzip.Middleware)
	a.router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.auth.RegisterHandler)
		r.Post("/login", a.auth.LoginHandler)
		r.Group(
			func(r chi.Router) {
				r.Use(a.auth.CheckMiddleware)
				r.Get("/ok", a.simpleHandler)
				r.Post("/orders", a.wallet.PostOrdersHandler)
				r.Get("/orders", a.wallet.GetOrdersHandler)
				r.Get("/balance", a.wallet.BalanceHandler)
				r.Post("/balance/withdraw", a.wallet.WithdrawHandler)
				r.Get("/withdrawals", a.wallet.GetWithdrawalsHandler)
			})
	})
}

func (a api) simpleHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (a api) SeverStart(ctx context.Context, apiURL string) error {
	defer a.logger.Sync()
	a.logger.Infof("Server started at %s", apiURL)
	httpServer := &http.Server{
		Addr:    apiURL,
		Handler: a.router,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := httpServer.ListenAndServe()
		if err != nil {
			return fmt.Errorf("serverapi.Run: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
	return nil
}
