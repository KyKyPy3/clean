//nolint:godot // file has comments for swagger doc
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

	"github.com/KyKyPy3/clean/internal/application/core"
	http_dto "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/game/application/command"
	"github.com/KyKyPy3/clean/internal/modules/game/application/query"
	"github.com/KyKyPy3/clean/internal/modules/game/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const (
	requestTimeout  = 10 * time.Second
	defaultPageSize = 50
)

type CommandBus interface {
	Dispatch(context.Context, core.Command) (any, error)
}

type QueryBus interface {
	Ask(context.Context, core.Query) (any, error)
}

type GameHandlers struct {
	Commands CommandBus
	Queries  QueryBus
	tracer   trace.Tracer
	Logger   logger.Logger
}

func NewGameHandlers(v1 *echo.Group, commands CommandBus, queries QueryBus, logger logger.Logger) {
	handlers := &GameHandlers{
		Commands: commands,
		Queries:  queries,
		Logger:   logger,
		tracer:   otel.Tracer(""),
	}

	v1.GET("/game", handlers.Fetch)
	v1.POST("/game", handlers.Create)
}

// Fetch godoc
// @Summary Fetch games
// @Description Fetch games handler
// @Tags Game
// @Accept json
// @Produce json
// @Success 201 {object} entity.Game
// @Router /game [get]
func (g *GameHandlers) Fetch(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	var errorList []*http_dto.ValidationError
	opts := dto.FetchGamesDTO{
		Limit:  defaultPageSize,
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
				errorList = append(errorList, &http_dto.ValidationError{
					Field:  bindingError.Field,
					Value:  bindingError.Values,
					Reason: "parse",
				})
			}
		}

		g.Logger.Errorf("failed to decode request params: %#v", errorList)

		return c.JSON(
			http.StatusBadRequest,
			http_dto.ResponseDTO{
				Status:  http.StatusBadRequest,
				Message: "error",
				Errors:  errorList,
			},
		)
	}

	err := c.Validate(opts)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
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

	q := query.FetchGamesQuery{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	}
	games, err := g.Queries.Ask(ctx, q)
	if err != nil {
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
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"games": games.([]dto.GameDTO),
			},
		},
	)
}

// Create godoc
// @Summary Create game
// @Description Create gamr handler
// @Tags Game
// @Accept json
// @Produce json
// @Success 201
// @Router /game [post]
// @BasePath /v1
func (g *GameHandlers) Create(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), requestTimeout)
	defer cancel()

	ctx, span := g.tracer.Start(ctx, "GameHandlers.Create")
	defer span.End()

	var errorList []*http_dto.ValidationError
	params := dto.CreateGameDTO{}

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
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
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

	g.Logger.Debugf("Create game with params %v", params)

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(
			http.StatusForbidden,
			http_dto.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	cmd := command.NewCreateGameCommand(params.Name, userID)
	_, err = g.Commands.Dispatch(ctx, cmd)
	if err != nil {
		g.Logger.Errorf("Failed to create registration %w", err)

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
