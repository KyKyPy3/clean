package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KyKyPy3/clean/config"
	"github.com/KyKyPy3/clean/internal/user/controller/http/dto"
	v1 "github.com/KyKyPy3/clean/internal/user/controller/http/v1"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	mocks "github.com/KyKyPy3/clean/mocks/internal_/user/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errInternal = errors.New("iternal error")
)

func TestFetchHandler(t *testing.T) {
	var mockUser entity.User
	mockUserList := make([]entity.User, 0)
	mockUserList = append(mockUserList, mockUser)

	// Create echo
	e := echo.New()

	// Create logger
	// TODO: add discard logger here
	loggerCfg := &config.LoggerConfig{Mode: "test"}
	logger := logger.NewLogger(loggerCfg)
	logger.Init()

	cases := []struct {
		name        string
		limit       string
		respStatus  int
		respMessage string
		respErrors  []*dto.ValidationError
		respError   string
		mockError   error
		mockResp    interface{}
	}{
		{
			name:        "Success",
			limit:       "10",
			respStatus:  http.StatusOK,
			respMessage: "success",
			respErrors:  nil,
			respError:   "",
			mockError:   nil,
			mockResp:    mockUserList,
		},
		{
			name:        "Failed",
			limit:       "10",
			respStatus:  http.StatusInternalServerError,
			respMessage: "error",
			respErrors:  nil,
			respError:   "iternal error",
			mockError:   errInternal,
			mockResp:    nil,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create user usecase mock
			userUsecaseMock := mocks.NewUserUsecase(t)

			handler := v1.UserHandlers{
				UserService: userUsecaseMock,
				Logger:      logger,
			}

			userUsecaseMock.On("Fetch", mock.Anything, mock.Anything).Return(tc.mockResp, tc.mockError).Once()

			req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/api/v1/user?limit="+tc.limit, strings.NewReader(""))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err = handler.Fetch(c)

			require.NoError(t, err)
			assert.Equal(t, tc.respStatus, rec.Code)

			var d *dto.ResponseDTO
			err = json.NewDecoder(rec.Body).Decode(&d)
			assert.NoError(t, err)

			assert.Equal(t, tc.respMessage, d.Message)
			assert.Equal(t, tc.respError, d.Error)

			userUsecaseMock.AssertExpectations(t)
		})
	}
}
