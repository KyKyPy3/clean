package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	common_http "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	mocks "github.com/KyKyPy3/clean/mocks/internal_/modules/user/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
)

var (
	errInternal = errors.New("internal error")
)

func TestFetchHandler(t *testing.T) {
	var mockUser entity.User
	mockUserList := make([]entity.User, 0)
	mockUserList = append(mockUserList, mockUser)

	// Create echo
	e := echo.New()

	// Create logger
	// TODO: add discard logger here
	log := logger.NewLogger(logger.Config{
		Mode: "test",
	})
	log.Init()

	cases := []struct {
		name        string
		limit       string
		offset      string
		respStatus  int
		respMessage string
		respErrors  []*common_http.ValidationError
		respError   string
		mockError   error
		mockResp    interface{}
	}{
		{
			name:        "Success",
			limit:       "10",
			offset:      "0",
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
			offset:      "0",
			respStatus:  http.StatusInternalServerError,
			respMessage: "error",
			respErrors:  nil,
			respError:   "internal error",
			mockError:   errInternal,
			mockResp:    nil,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create user command bus mock
			userCommandBusMock := mocks.NewCommandBus(t)
			userQueryBusMock := mocks.NewQueryBus(t)

			handler := v1.UserHandlers{
				Commands: userCommandBusMock,
				Queries:  userQueryBusMock,
				Logger:   log,
			}

			userQueryBusMock.On("Ask", mock.Anything, mock.Anything).Return(tc.mockResp, tc.mockError).Once()

			req, err := http.NewRequestWithContext(
				context.TODO(),
				echo.GET,
				fmt.Sprintf("/api/v1/user?limit=%s&offset=%s", tc.limit, tc.offset),
				strings.NewReader(""),
			)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err = handler.Fetch(c)

			require.NoError(t, err)
			assert.Equal(t, tc.respStatus, rec.Code)

			var d *common_http.ResponseDTO
			err = json.NewDecoder(rec.Body).Decode(&d)
			assert.NoError(t, err)

			assert.Equal(t, tc.respMessage, d.Message)
			assert.Equal(t, tc.respError, d.Error)

			userQueryBusMock.AssertExpectations(t)
		})
	}
}
