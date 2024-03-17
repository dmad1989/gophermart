package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/jsonobject"
	"github.com/dmad1989/gophermart/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	okReqBody = `{"login":"hit", "password":"hit"}`
)

type postRequest struct {
	body       io.Reader
	jsonHeader bool
}

type createUserMockParams struct {
	resID     int
	resErr    error
	callTimes int
}
type getUserMockParams struct {
	resUser   jsonobject.User
	resErr    error
	callTimes int
}

type expectedPostResponse struct {
	code         int
	errorMessage string
}

func TestRegisterHandler(t *testing.T) {
	ctx := initContext()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mDB := mocks.NewMockDB(ctrl)
	tests := []struct {
		name       string
		request    postRequest
		expResp    expectedPostResponse
		mockParams createUserMockParams
	}{
		{
			name: "negative - no json header",
			request: postRequest{
				jsonHeader: false,
				body:       strings.NewReader("JSONBodyRequest"),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: ErrorRequestContentType.Error(),
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - emptyBody",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(""),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: "decoding request: EOF",
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - wrong json",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader("{?}"),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: "decoding request: parse error: syntax error near offset 1 of '{?}'",
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - json login empty",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(`{"password":"hit"}`),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: ErrorRequestLogin.Error(),
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - json password empty",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(`{"login":"hit"}`),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: ErrorRequestPassword.Error(),
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - db error pgconn.PgError UniqueViolation",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusConflict,
				errorMessage: ErrorUserLoginUnique.Error(),
			},
			mockParams: createUserMockParams{
				resID: 0,
				resErr: &pgconn.PgError{
					Code: pgerrcode.UniqueViolation,
				},
				callTimes: 1,
			},
		},
		{
			name: "negative - db custom error",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusInternalServerError,
				errorMessage: "register: custom error",
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    errors.New("custom error"),
				callTimes: 1,
			},
		},
		{
			name: "negative - token generate error",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusInternalServerError,
				errorMessage: "generating token: generateToken: user id is 0",
			},
			mockParams: createUserMockParams{
				resID:     0,
				resErr:    nil,
				callTimes: 1,
			},
		},
		{
			name: "positive",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusOK,
				errorMessage: "",
			},
			mockParams: createUserMockParams{
				resID:     1,
				resErr:    nil,
				callTimes: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", tt.request.body)
			if tt.request.jsonHeader {
				req.Header.Set("Content-Type", "application/json")
			}
			mDB.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(tt.mockParams.resID, tt.mockParams.resErr).
				MaxTimes(tt.mockParams.callTimes)
			a := New(ctx, mDB)

			w := httptest.NewRecorder()
			a.RegisterHandler(w, req)
			res := w.Result()

			assert.Equal(t, tt.expResp.code, res.StatusCode, "statusCode mismatch")
			b, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			resBody := string(b)
			err = res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.expResp.errorMessage, string(resBody))
			if res.StatusCode == http.StatusOK {
				tokenInCookie := false
				for _, c := range res.Cookies() {
					if c.Name == "token" {
						tokenInCookie = true
						break
					}
				}
				assert.True(t, tokenInCookie, "no cookie in response")
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	ctx := initContext()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mDB := mocks.NewMockDB(ctrl)

	pass := sha256.Sum256([]byte("hit"))

	tests := []struct {
		name       string
		request    postRequest
		expResp    expectedPostResponse
		mockParams getUserMockParams
	}{
		{
			name: "negative - no json header",
			request: postRequest{
				jsonHeader: false,
				body:       strings.NewReader("JSONBodyRequest"),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: ErrorRequestContentType.Error(),
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{},
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - emptyBody",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(""),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: "decoding request: EOF",
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{},
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - wrong json",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader("{?}"),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: "decoding request: parse error: syntax error near offset 1 of '{?}'",
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{},
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - json login empty",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(`{"password":"hit"}`),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: ErrorRequestLogin.Error(),
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{},
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - json password empty",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(`{"login":"hit"}`),
			},
			expResp: expectedPostResponse{
				code:         http.StatusBadRequest,
				errorMessage: ErrorRequestPassword.Error(),
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{},
				resErr:    nil,
				callTimes: 0,
			},
		},
		{
			name: "negative - DB custom error",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusInternalServerError,
				errorMessage: "check user in DB: custom error",
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{},
				resErr:    errors.New("custom error"),
				callTimes: 1,
			},
		},
		{
			name: "negative - wrong pass",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusUnauthorized,
				errorMessage: ErrorUserPassword.Error(),
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{Login: "wip", HashPassword: []byte{}},
				resErr:    nil,
				callTimes: 1,
			},
		},
		{
			name: "negative - generate token",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusInternalServerError,
				errorMessage: "generating token: generateToken: user id is 0",
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{ID: 0, Login: "hit", HashPassword: pass[:]},
				resErr:    nil,
				callTimes: 1,
			},
		},
		{
			name: "positive",
			request: postRequest{
				jsonHeader: true,
				body:       strings.NewReader(okReqBody),
			},
			expResp: expectedPostResponse{
				code:         http.StatusOK,
				errorMessage: "",
			},
			mockParams: getUserMockParams{
				resUser:   jsonobject.User{ID: 1, Login: "hit", HashPassword: pass[:]},
				resErr:    nil,
				callTimes: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", tt.request.body)
			if tt.request.jsonHeader {
				req.Header.Set("Content-Type", "application/json")
			}
			mDB.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(tt.mockParams.resUser, tt.mockParams.resErr).MaxTimes(tt.mockParams.callTimes)
			a := New(ctx, mDB)

			w := httptest.NewRecorder()
			a.LoginHandler(w, req)
			res := w.Result()

			assert.Equal(t, tt.expResp.code, res.StatusCode, "statusCode mismatch")
			b, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			resBody := string(b)
			err = res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.expResp.errorMessage, string(resBody))

			if res.StatusCode == http.StatusOK {
				tokenInCookie := false
				for _, c := range res.Cookies() {
					if c.Name == "token" {
						tokenInCookie = true
						break
					}
				}
				assert.True(t, tokenInCookie, "no cookie in response")
			}
		})
	}
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
