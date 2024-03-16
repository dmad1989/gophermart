package auth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

var (
	ErrorNoUser             = errors.New("no userid in auth token")
	ErrorInvalidToken       = errors.New("auth token not valid")
	ErrorUserLoginUnique    = errors.New("user login is not unique")
	ErrorUserPassword       = errors.New("wrong password")
	ErrorRequestContentType = errors.New("content-type have to be application/json")
	ErrorRequestLogin       = errors.New("user login in request is empty")
	ErrorRequestPassword    = errors.New("user password in request is empty")
)

const tokenExp = time.Hour * 6
const secretKey = "gopracticgphermarksecretkey"

type auth struct {
	db     DB
	logger *zap.SugaredLogger
}

type DB interface {
	CreateUser(ctx context.Context, user jsonobject.User) (int, error)
	GetUserByLogin(ctx context.Context, login string) (jsonobject.User, error)
}

func New(ctx context.Context, db DB) *auth {
	return &auth{db: db, logger: ctx.Value(config.LoggerCtxKey).(*zap.SugaredLogger)}
}

func (a auth) RegisterHandler(res http.ResponseWriter, req *http.Request) {
	user, err := getUserFromRequest(req)
	if err != nil {
		errorResponse(res, http.StatusBadRequest, err)
		return
	}
	// шифруем пароль
	hashed := sha256.Sum256([]byte(user.Password))
	user.HashPassword = hashed[:]
	// записываем в БД
	user.ID, err = a.db.CreateUser(req.Context(), user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			errorResponse(res, http.StatusConflict, ErrorUserLoginUnique)
			return
		}
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("register: %w", err))
		return
	}
	// генерируем токен
	token, err := generateToken(user.ID)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("generating token: %w", err))
		return
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
}

func (a auth) LoginHandler(res http.ResponseWriter, req *http.Request) {
	user, err := getUserFromRequest(req)
	if err != nil {
		errorResponse(res, http.StatusBadRequest, fmt.Errorf("check user in DB: %w", err))
		return
	}
	// проверяем пользователя
	userDB, err := a.db.GetUserByLogin(req.Context(), user.Login)
	if err != nil {
		errorResponse(res, http.StatusUnauthorized, fmt.Errorf("check user in DB: %w", err))
		return
	}
	// проверяем пароль
	hashed := sha256.Sum256([]byte(user.Password))

	if !bytes.Equal(hashed[:], userDB.HashPassword) {
		errorResponse(res, http.StatusUnauthorized, ErrorUserPassword)
		return
	}
	// генерируем токен
	token, err := generateToken(userDB.ID)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, fmt.Errorf("generating token:  %w", err))
		return
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
}

func (a auth) CheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextW := w
		//проверяем токен
		h.ServeHTTP(nextW, r.WithContext(r.Context()))
	})
}

func getUserFromRequest(req *http.Request) (jsonobject.User, error) {
	var user jsonobject.User
	if req.Header.Get("Content-Type") != "application/json" {
		return user, ErrorRequestContentType
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return user, fmt.Errorf("reading request body: %w", err)
	}
	if err := user.UnmarshalJSON(body); err != nil {
		return user, fmt.Errorf("cutterJsonHandler: decoding request: %w", err)
	}
	if user.Login == "" {
		return user, ErrorRequestLogin
	}
	if user.Password == "" {
		return user, ErrorRequestPassword
	}
	// todo check not empty?
	return user, nil
}

func generateToken(userID int) (string, error) {
	if userID == 0 {
		return "", errors.New("generateToken: user id is 0")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("generateToken: %w", err)
	}
	return tokenString, nil
}

func errorResponse(res http.ResponseWriter, status int, err error) {
	res.WriteHeader(status)
	res.Write([]byte(err.Error()))
}
