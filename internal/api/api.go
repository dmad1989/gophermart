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
type Auth interface {
	LoginHandler(http.ResponseWriter, *http.Request)
	RegisterHandler(http.ResponseWriter, *http.Request)
	CheckMiddleware(h http.Handler) http.Handler
}

type api struct {
	router *chi.Mux

	app        App
	auth       Auth
	accrualURL string
}

func New(ctx context.Context, app App, accrualURL string, auth Auth) *api {
	api := &api{
		app:        app,
		router:     chi.NewRouter(),
		auth:       auth,
		accrualURL: accrualURL}
	api.initRouter()
	return api
}

func (a api) initRouter() {
	a.router.Use(middleware.Logger, middleware.Recoverer) // todo auth, gzip
	// a.router.Get("/ok", a.simpleHandler)
	a.router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.auth.RegisterHandler)
		r.Post("/login", a.auth.LoginHandler)
		r.Group(
			func(r chi.Router) {
				r.Use(a.auth.CheckMiddleware)
				r.Get("/ok", a.simpleHandler)
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
