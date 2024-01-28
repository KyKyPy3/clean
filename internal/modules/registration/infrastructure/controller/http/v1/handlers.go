package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	http_dto "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/registration/application"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/command"
	"github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type RegistrationApp interface {
	Create(ctx context.Context, email string) error
}

type RegistrationHandlers struct {
	registrationApp *application.Application
	logger          logger.Logger
	tracer          trace.Tracer
}

func NewRegistrationHandlers(v1 *echo.Group, registrationApp *application.Application, logger logger.Logger) {
	handlers := &RegistrationHandlers{
		registrationApp: registrationApp,
		logger:          logger,
		tracer:          otel.Tracer(""),
	}

	v1.POST("/registration", handlers.Create)
	v1.GET("/registration/:id", handlers.Confirm)
}

// Create godoc
// @Summary Create registration
// @Description Create registration handler
// @Tags Registration
// @Accept json
// @Produce json
// @Success 201
// @Router /registration [post]
func (r *RegistrationHandlers) Create(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	ctx, span := r.tracer.Start(ctx, "RegistrationHandlers.Create")
	defer span.End()

	var errorList []*http_dto.ValidationError
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
			http_dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Error:   validationErr,
			},
		)
	}

	err = c.Validate(params)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorList = append(errorList, &http_dto.ValidationError{
				Field:  e.Field(),
				Value:  e.Value(),
				Reason: e.Tag(),
			})
		}

		return c.JSON(
			http.StatusBadRequest,
			http_dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	r.logger.Debugf("Create registration with params %v", params)

	cmd := command.CreateRegistrationCommand{
		Email: params.Email,
	}
	err = r.registrationApp.Commands.CreateRegistration.Handle(ctx, cmd)
	if err != nil {
		r.logger.Errorf("Failed to create registration %w", err)

		return c.JSON(
			http.StatusInternalServerError,
			http_dto.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		http_dto.ResponseDTO{
			Status:  http.StatusCreated,
			Message: "success",
		},
	)
}

// Confirm godoc
// @Summary Confirm registration
// @Description Confirm registration handler
// @Tags Registration
// @Accept json
// @Produce json
// @Success 201
// @Router /registration [get]
func (r *RegistrationHandlers) Confirm(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	ctx, span := r.tracer.Start(ctx, "RegistrationHandlers.Сonfirm")
	defer span.End()

	id := c.Param("id")

	r.logger.Debugf("Confirm registration with id '%v'", id)

	cmd := command.ConfirmRegistrationCommand{
		ID: id,
	}
	err := r.registrationApp.Commands.ConfirmRegistration.Handle(ctx, cmd)
	if err != nil {
		r.logger.Errorf("Failed to confirm registration %w", err)

		return c.JSON(
			http.StatusInternalServerError,
			http_dto.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Error:   err.Error(),
			},
		)
	}

	return c.NoContent(http.StatusOK)
}
