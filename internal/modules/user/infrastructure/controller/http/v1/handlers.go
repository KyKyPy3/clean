package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/KyKyPy3/clean/internal/modules/user/application/query"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	common_http "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/user/application/command"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const requestTimeout = 10 * time.Second

type CommandBus interface {
	Dispatch(context.Context, core.Command) (any, error)
}

type QueryBus interface {
	Ask(context.Context, core.Query) (any, error)
}

type UserHandlers struct {
	Commands CommandBus
	Queries  QueryBus
	Logger   logger.Logger
}

func NewUserHandlers(v1 *echo.Group, commands CommandBus, queries QueryBus, logger logger.Logger) {
	handlers := &UserHandlers{Commands: commands, Queries: queries, Logger: logger}

	v1.GET("/user", handlers.Fetch)
	v1.POST("/user/:id", handlers.Update)
	v1.GET("/user/:id", handlers.GetByID)
	v1.DELETE("/user/:id", handlers.Delete)
}

// Fetch godoc
// @Summary Fetch users
// @Description Fetch users handler
// @Tags User
// @Accept json
// @Produce json
// @Success 201 {object} entity.User
// @Router /user [get]
func (h *UserHandlers) Fetch(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

	q := query.FetchUsersQuery{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	}
	users, err := h.Queries.Ask(ctx, q)
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

	for _, user := range users.([]entity.User) {
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

// Update godoc
// @Summary Update user
// @Description Update user handler
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Router /user/{id} [post]
func (h *UserHandlers) Update(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	var errorList []*common_http.ValidationError
	params := dto.UpdateUserDTO{}

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

	h.Logger.Debugf("Update user with params %v", params)

	cmd := command.UpdateUserCommand{
		Name:       params.Name,
		Surname:    params.Surname,
		Middlename: params.Middlename,
		Email:      params.Email,
	}

	_, err = h.Commands.Dispatch(ctx, cmd)
	if err != nil {
		h.Logger.Errorf("Failed to update user %w", err)

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
			Status:  http.StatusOK,
			Message: "success",
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
