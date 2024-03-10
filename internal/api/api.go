package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

type App interface{}

type api struct {
	app        App
	accrualURL string
	router     *chi.Mux
}

func New(ctx context.Context, app App, accrualURL string) *api {
	api := &api{
		app:        app,
		router:     chi.NewRouter(),
		accrualURL: accrualURL}
	api.initRouter()
	return api
}

func (a api) initRouter() {
	a.router.Use(middleware.Logger, middleware.Recoverer) // todo auth, gzip
	a.router.Get("/ok", a.simpleHandler)
	a.router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.registerHandler)
		r.Post("/login", a.authHandler)
		r.Group(
			func(r chi.Router) {
				r.Use(a.authCheckMiddleware)
				r.Post("/orders", a.postOrdersHandler)
				r.Get("/orders", a.getOrdersHandler)
				r.Get("/balance", a.balanceHandler)
				r.Post("/balance/withdraw", a.withdrawHandler)
				r.Get("/withdrawals", a.allWithdrawalsHandler)
			})
	})

}

func (a api) simpleHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (a api) SeverStart(ctx context.Context, apiURL string) error {
	// defer logging.Log.Sync()
	// logging.Log.Infof("Server started at %s", apiURL)
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

func (a api) postOrdersHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (a api) getOrdersHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (a api) balanceHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (a api) withdrawHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (a api) allWithdrawalsHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}
