package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	Fetch(ctx context.Context, limit int64) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id common.ID) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Delete(ctx context.Context, id common.ID) error
}

type handlers struct {
	userService UserUsecase
}

func NewUserHandlers(v1 *echo.Group, userService UserUsecase) {
	handlers := &handlers{userService: userService}

	v1.GET("/user", handlers.Fetch)
	v1.POST("/user", handlers.Create)
	v1.GET("/user/:id", handlers.Get)
	v1.DELETE("/user/:id", handlers.Delete)
}

func (h *handlers) Fetch(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	users, err := h.userService.Fetch(ctx, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *handlers) Create(c echo.Context) error {
	return c.JSON(http.StatusCreated, "")
}

func (h *handlers) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func (h *handlers) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
