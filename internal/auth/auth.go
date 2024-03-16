package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmad1989/gophermart/internal/jsonobject"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	argonpass "github.com/dwin/goArgonPass"
)

type auth struct {
	db DB
}

type DB interface {
	CreateUser(ctx context.Context, user jsonobject.User) error
}

func New(ctx context.Context, db DB) *auth {
	return &auth{db: db}
}

func (a auth) RegisterHandler(res http.ResponseWriter, req *http.Request) {
	user, err := getUserFromRequest(req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}
	// шифруем пароль
	//TODO отказаться от argonpass, использовать SHA256
	user.HashPassword, err = argonpass.Hash(user.Password, nil)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Errorf("argonpass, hashing password:  %w", err).Error()))
		return
	}
	// записываем в БД
	err = a.db.CreateUser(req.Context(), user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte(errors.New("user login is not unique!").Error()))
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Errorf("register:  %w", err).Error()))
		return
	}
	// генерируем токен

	res.WriteHeader(http.StatusOK)
}

func (a auth) LoginHandler(res http.ResponseWriter, req *http.Request) {
	user, err := getUserFromRequest(req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}
	// проверяем пользователя
	// проверяем пароль
	//TODO отказаться от argonpass, использовать SHA256
	err = argonpass.Verify(user.Password, hash)
	if err != nil {
		if errors.Is(err, argonpass.ErrHashMismatch) {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte(fmt.Errorf("wrong password:  %w", err).Error()))
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Errorf("argonpass, checking password:  %w", err).Error()))
		return
	}
	// генерируем токен
	res.WriteHeader(http.StatusOK)
}

func getUserFromRequest(req *http.Request) (jsonobject.User, error) {
	var user jsonobject.User
	if req.Header.Get("Content-Type") != "application/json" {
		return user, errors.New("content-type have to be application/json")
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return user, fmt.Errorf("reading request body: %w", err)
	}
	if err := user.UnmarshalJSON(body); err != nil {
		return user, fmt.Errorf("cutterJsonHandler: decoding request: %w", err)
	}
	return user, nil
}

func (a auth) CheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextW := w
		//проверяем токен
		h.ServeHTTP(nextW, r.WithContext(r.Context()))
	})
}
