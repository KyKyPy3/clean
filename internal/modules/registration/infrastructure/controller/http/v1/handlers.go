package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	common_dto "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type RegistrationUsecase interface {
	Create(ctx context.Context, user entity.Registration) error
}

type RegistrationHandlers struct {
	RegistrationService RegistrationUsecase
	Logger              logger.Logger
}

func NewRegistrationHandlers(v1 *echo.Group, registrationService RegistrationUsecase, logger logger.Logger) {
	handlers := &RegistrationHandlers{RegistrationService: registrationService, Logger: logger}

	v1.POST("/registration", handlers.Create)
}

// Create godoc
// @Summary Create registration
// @Description Create registration handler
// @Tags Registration
// @Accept json
// @Produce json
// @Success 201 {object} entity.Registration
// @Router /registration [post]
func (r *RegistrationHandlers) Create(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var errorList []*common_dto.ValidationError
	params := dto.CreateRegistrationDTO{}

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
			common_dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Error:   validationErr,
			},
		)
	}

	err = c.Validate(params)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorList = append(errorList, &common_dto.ValidationError{
				Field:  e.Field(),
				Value:  e.Value(),
				Reason: e.Tag(),
			})
		}

		return c.JSON(
			http.StatusBadRequest,
			common_dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	registration, err := dto.RegistrationFromRequest(params)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			common_dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	r.Logger.Debugf("Create registration with params %v", registration)

	err = r.RegistrationService.Create(ctx, registration)
	if err != nil {
		r.Logger.Errorf("Failed to create registration %w", err)

		return c.JSON(
			http.StatusInternalServerError,
			common_dto.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		common_dto.ResponseDTO{
			Status:  http.StatusCreated,
			Message: "success",
			Data: map[string]interface{}{
				"user": dto.RegistrationToResponse(registration),
			},
		},
	)
}
