package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/session/application/ports"
	"github.com/KyKyPy3/clean/pkg/jwt"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type AuthMiddleware struct {
	sessionStorage ports.SessionRedisStorage
	logger         logger.Logger
	jwt            *jwt.JWT
}

func NewAuthMiddleware(jwt *jwt.JWT, sessionStorage ports.SessionRedisStorage, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:            jwt,
		sessionStorage: sessionStorage,
		logger:         logger,
	}
}

func (a *AuthMiddleware) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var accessToken string

		header := c.Request().Header
		authorization := header.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			accessToken = strings.TrimPrefix(authorization, "Bearer ")
		} else {
			cookie, err := c.Cookie("access_token")

			if err != nil {
				return c.NoContent(
					http.StatusUnauthorized,
				)
			}
			accessToken = cookie.Value
		}

		if accessToken == "" {
			return c.NoContent(
				http.StatusUnauthorized,
			)
		}

		token, err := a.jwt.ValidateToken(accessToken)
		if err != nil {
			return c.NoContent(
				http.StatusForbidden,
			)
		}

		tokenID, err := common.ParseUID(token.TokenUUID)
		if err != nil {
			return c.NoContent(
				http.StatusForbidden,
			)
		}

		t, err := a.sessionStorage.Get(c.Request().Context(), tokenID)
		if err != nil {
			a.logger.Debugf("%s", err)

			return c.NoContent(
				http.StatusForbidden,
			)
		}

		c.Set("user_id", t.UserID().String())
		c.Set("access_token_id", t.ID().String())

		return next(c)
	}
}
