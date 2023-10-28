package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/controller/http/dto"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	Fetch(ctx context.Context, limit int64) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id common.ID) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Delete(ctx context.Context, id common.ID) error
}

type UserHandlers struct {
	UserService UserUsecase
	Logger      logger.Logger
}

func NewUserHandlers(v1 *echo.Group, userService UserUsecase, logger logger.Logger) {
	handlers := &UserHandlers{UserService: userService, Logger: logger}

	v1.GET("/user", handlers.Fetch)
	v1.POST("/user", handlers.Create)
	v1.GET("/user/:id", handlers.GetByID)
	v1.DELETE("/user/:id", handlers.Delete)
}

func (h *UserHandlers) Fetch(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var errorList []*dto.ValidationError
	opts := dto.FetchUsersDTO{
		Limit: 50,
	}

	// Parse given params
	errs := echo.QueryParamsBinder(c).
		FailFast(false).
		Int64("limit", &opts.Limit).
		BindErrors()
	if errs != nil {
		for _, err := range errs {
			var bindingError *echo.BindingError
			if errors.As(err, &bindingError) {
				errorList = append(errorList, &dto.ValidationError{
					Field:  bindingError.Field,
					Value:  bindingError.Values,
					Reason: "parse",
				})
			}
		}

		h.Logger.Errorf("failed to decode request params: %#v", errorList)

		return c.JSON(
			http.StatusBadRequest,
			dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	err := c.Validate(opts)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorList = append(errorList, &dto.ValidationError{
				Field:  e.Field(),
				Value:  e.Value(),
				Reason: e.Tag(),
			})
		}

		return c.JSON(
			http.StatusBadRequest,
			dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	users, err := h.UserService.Fetch(ctx, opts.Limit)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.ResponseDTO{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"users": users,
			},
		},
	)
}

func (h *UserHandlers) Create(c echo.Context) error {
	return c.JSON(http.StatusCreated, "")
}

func (h *UserHandlers) GetByID(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func (h *UserHandlers) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
