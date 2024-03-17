package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmad1989/gophermart/internal/config"
	"github.com/dmad1989/gophermart/internal/mocks"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", tt.request.body)
			if tt.request.jsonHeader {
				req.Header.Set("Content-Type", "application/json")
			}
			mDB.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.mockParams.resID, tt.mockParams.resErr).MaxTimes(tt.mockParams.callTimes)
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
