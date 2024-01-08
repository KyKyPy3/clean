package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/domain/common"
	common_http "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type UserUsecase interface {
	Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id common.UID) (entity.User, error)
	GetByEmail(ctx context.Context, email common.Email) (entity.User, error)
	Delete(ctx context.Context, id common.UID) error
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

	var errorList []*common_http.ValidationError
	opts := dto.FetchUsersDTO{
		Limit:  50,
		Offset: 0,
	}

	// Parse given params
	errs := echo.QueryParamsBinder(c).
		FailFast(false).
		Int64("limit", &opts.Limit).
		Int64("offset", &opts.Offset).
		BindErrors()
	if errs != nil {
		for _, err := range errs {
			var bindingError *echo.BindingError
			if errors.As(err, &bindingError) {
				errorList = append(errorList, &common_http.ValidationError{
					Field:  bindingError.Field,
					Value:  bindingError.Values,
					Reason: "parse",
				})
			}
		}

		h.Logger.Errorf("failed to decode request params: %#v", errorList)

		return c.JSON(
			http.StatusBadRequest,
			common_http.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	err := c.Validate(opts)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorList = append(errorList, &common_http.ValidationError{
				Field:  e.Field(),
				Value:  e.Value(),
				Reason: e.Tag(),
			})
		}

		return c.JSON(
			http.StatusBadRequest,
			common_http.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	users, err := h.UserService.Fetch(ctx, opts.Limit, opts.Offset)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			common_http.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Error:   err.Error(),
			},
		)
	}

	respUsers := make([]dto.UserDTO, 0)

	for _, user := range users {
		respUsers = append(respUsers, dto.UserToResponse(user))
	}

	return c.JSON(
		http.StatusOK,
		common_http.ResponseDTO{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"users": respUsers,
			},
		},
	)
}

// Create godoc
// @Summary Create user
// @Description Create user handler
// @Tags User
// @Accept json
// @Produce json
// @Success 201 {object} entity.User
// @Router /user [post]
func (h *UserHandlers) Create(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var errorList []*common_http.ValidationError
	params := dto.CreateUserDTO{}

	// Parse given params
	err := c.Bind(&params)
	if err != nil {
		var bindingError *echo.HTTPError
		var validationErr string
		if errors.As(err, &bindingError) {
			validationErr = fmt.Sprint(bindingError.Message)
		} else {
			validationErr = err.Error()
		}

		return c.JSON(
			http.StatusBadRequest,
			common_http.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Error:   validationErr,
			},
		)
	}

	err = c.Validate(params)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorList = append(errorList, &common_http.ValidationError{
				Field:  e.Field(),
				Value:  e.Value(),
				Reason: e.Tag(),
			})
		}

		return c.JSON(
			http.StatusBadRequest,
			common_http.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	user, err := dto.UserFromRequest(params)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			common_http.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	h.Logger.Debugf("Create user with params %v", user)

	user, err = h.UserService.Create(ctx, user)
	if err != nil {
		h.Logger.Errorf("Failed to create user %w", err)

		return c.JSON(
			http.StatusInternalServerError,
			common_http.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		common_http.ResponseDTO{
			Status:  http.StatusCreated,
			Message: "success",
			Data: map[string]interface{}{
				"user": dto.UserToResponse(user),
			},
		},
	)
}

// GetByID godoc
// @Summary Get by id user
// @Description Get by id user handler
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "user_id"
// @Success 200 {object} entity.User
// @Router /user/{id} [get]
func (h *UserHandlers) GetByID(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func (h *UserHandlers) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
