package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	common_http "github.com/KyKyPy3/clean/internal/infrastructure/controller/http"
	"github.com/KyKyPy3/clean/internal/modules/session/application/command"
	"github.com/KyKyPy3/clean/internal/modules/session/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/jwt"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const (
	requestTimeout   = 10 * time.Second
	cookieExpiration = -time.Hour * 24
	accessTokenKey   = "access_token"
	refreshTokenKey  = "refresh_token"
)

type CommandBus interface {
	Dispatch(context.Context, core.Command) (any, error)
}

type QueryBus interface {
	Ask(context.Context, core.Query) (any, error)
}

type AuthHandlers struct {
	Cfg      *config.Config
	Commands CommandBus
	Queries  QueryBus
	Jwt      *jwt.JWT
	Logger   logger.Logger
}

func NewAuthHandlers(
	publicMountPoint *echo.Group,
	privateMountPoint *echo.Group,
	commands CommandBus,
	queries QueryBus,
	cfg *config.Config,
	jwt *jwt.JWT,
	logger logger.Logger,
) {
	handlers := &AuthHandlers{Commands: commands, Queries: queries, Cfg: cfg, Jwt: jwt, Logger: logger}

	publicMountPoint.POST("/auth/login", handlers.Login)
	privateMountPoint.POST("/auth/logout", handlers.Logout)
	privateMountPoint.POST("/auth/refresh", handlers.RefreshToken)
}

func (a *AuthHandlers) Logout(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	cookie, err := c.Cookie(refreshTokenKey)
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			common_http.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	token, err := a.Jwt.ValidateToken(cookie.Value)
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			common_http.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	// Remove token from database
	accessTokenID, ok := c.Get("access_token_id").(string)
	if !ok {
		return c.JSON(
			http.StatusInternalServerError,
			common_http.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
			},
		)
	}

	cmd := command.LogoutUserCommand{
		RefreshTokenID: token.TokenUUID,
		AccessTokenID:  accessTokenID,
	}
	_, err = a.Commands.Dispatch(ctx, cmd)
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

	// Delete cookie
	expired := time.Now().Add(cookieExpiration)

	cookie = new(http.Cookie)
	cookie.Name = accessTokenKey
	cookie.Value = ""
	cookie.Expires = expired
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = refreshTokenKey
	cookie.Value = ""
	cookie.Expires = expired
	c.SetCookie(cookie)

	return c.JSON(
		http.StatusOK,
		common_http.ResponseDTO{
			Status:  http.StatusOK,
			Message: "success",
		},
	)
}

func (a *AuthHandlers) RefreshToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	cookie, err := c.Cookie(refreshTokenKey)
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			common_http.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	token, err := a.Jwt.ValidateToken(cookie.Value)
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			common_http.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	// Refresh access token
	cmd := command.RefreshSessionCommand{
		ID:         token.TokenUUID,
		UserID:     token.UserID,
		AccessTTL:  a.Cfg.Jwt.AccessTokenMaxAge,
		RefreshTTL: a.Cfg.Jwt.RefreshTokenMaxAge,
	}
	res, err := a.Commands.Dispatch(ctx, cmd)
	if err != nil {
		a.Logger.Errorf("Failed to refresh token %v", err)

		return c.JSON(
			http.StatusInternalServerError,
			common_http.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	// Unpack results
	meta, ok := res.(command.RefreshSessionResult)
	if !ok {
		return c.JSON(
			http.StatusInternalServerError,
			common_http.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
			},
		)
	}

	// Generate tokens
	accessToken, err := a.Jwt.CreateToken(
		meta.AccessToken.ID().String(),
		meta.UserID.String(),
		a.Cfg.Jwt.AccessTokenMaxAge,
	)
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

	// Set tokens to cookie
	cookie = new(http.Cookie)
	cookie.Name = refreshTokenKey
	cookie.Value = *accessToken.Token
	cookie.Path = "/"
	cookie.MaxAge = int(a.Cfg.Jwt.AccessTokenMaxAge.Seconds())
	c.SetCookie(cookie)

	return c.JSON(
		http.StatusOK,
		common_http.ResponseDTO{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"access_token": accessToken.Token,
			},
		},
	)
}

// Login godoc
// @Summary Login user
// @Description Login user handler
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.Token
// @Router /auth/login [post]
// TOFIX: refactor function
//
//nolint:funlen // can't make it simpler now
func (a *AuthHandlers) Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	params := dto.LoginDTO{}

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

	if err = c.Validate(params); err != nil {
		return handleValidationErrors(c, err)
	}

	a.Logger.Debugf("Login with params %v", params)

	cmd := command.LoginUserCommand{
		Email:      params.Email,
		Password:   params.Password,
		AccessTTL:  a.Cfg.Jwt.AccessTokenMaxAge,
		RefreshTTL: a.Cfg.Jwt.RefreshTokenMaxAge,
	}
	res, err := a.Commands.Dispatch(ctx, cmd)
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			common_http.ResponseDTO{
				Status:  http.StatusForbidden,
				Message: "error",
			},
		)
	}

	// Generate tokens
	meta, ok := res.(command.LoginUserResult)
	if !ok {
		return c.JSON(
			http.StatusInternalServerError,
			common_http.ResponseDTO{
				Status:  http.StatusInternalServerError,
				Message: "error",
			},
		)
	}

	accessToken, err := a.Jwt.CreateToken(
		meta.AccessToken.ID().String(),
		meta.UserID.String(),
		a.Cfg.Jwt.AccessTokenMaxAge,
	)
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

	refreshToken, err := a.Jwt.CreateToken(
		meta.RefreshToken.ID().String(),
		meta.UserID.String(),
		a.Cfg.Jwt.RefreshTokenMaxAge,
	)
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

	a.setCookie(c, accessToken, refreshToken)

	return c.JSON(
		http.StatusOK,
		common_http.ResponseDTO{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"access_token":  accessToken.Token,
				"refresh_token": refreshToken.Token,
			},
		},
	)
}

func (a *AuthHandlers) setCookie(c echo.Context, accessToken, refreshToken *jwt.Token) {
	cookie := new(http.Cookie)
	cookie.Name = accessTokenKey
	cookie.Value = *accessToken.Token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = int(a.Cfg.Jwt.AccessTokenMaxAge.Seconds())
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = refreshTokenKey
	cookie.Value = *refreshToken.Token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = int(a.Cfg.Jwt.RefreshTokenMaxAge.Seconds())
	c.SetCookie(cookie)
}

func handleValidationErrors(c echo.Context, err error) error {
	status := http.StatusBadRequest
	message := "error"
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorList []*common_http.ValidationError
		for _, e := range validationErrors {
			errorList = append(errorList, &common_http.ValidationError{
				Field:  e.Field(),
				Value:  e.Value(),
				Reason: e.Tag(),
			})
		}
		return c.JSON(status, common_http.ResponseDTO{
			Status:  status,
			Message: message,
			Errors:  errorList,
		})
	}
	return c.JSON(status, common_http.ResponseDTO{
		Status:  status,
		Message: message,
		Error:   err.Error(),
	})
}
