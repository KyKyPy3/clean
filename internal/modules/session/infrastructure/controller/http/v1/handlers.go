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
}

func (a *AuthHandlers) Logout(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cookie, err := c.Cookie("refresh_token")
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
				Error:   err.Error(),
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
	expired := time.Now().Add(-time.Hour * 24)

	cookie = new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = ""
	cookie.Expires = expired
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = "refresh_token"
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
	//cookie, err := c.Cookie("refresh_token")
	//if err != nil {
	//	return c.JSON(
	//		http.StatusForbidden,
	//		common_http.ResponseDTO{
	//			Status:  http.StatusForbidden,
	//			Message: "error",
	//		},
	//	)
	//}
	//
	//token, err := a.Jwt.ValidateToken(cookie.Value)
	//if err != nil {
	//	return c.JSON(
	//		http.StatusForbidden,
	//		common_http.ResponseDTO{
	//			Status:  http.StatusForbidden,
	//			Message: "error",
	//		},
	//	)
	//}
	//
	//// Get token from database
	//
	//// Find user from database
	//
	//// Generate tokens
	//accessToken, err := a.Jwt.CreateToken(user.(*entity.User).ID().String(), a.Cfg.Jwt.AccessTokenMaxAge)
	//if err != nil {
	//	return c.JSON(
	//		http.StatusInternalServerError,
	//		common_http.ResponseDTO{
	//			Status:  http.StatusInternalServerError,
	//			Message: "error",
	//			Error:   err.Error(),
	//		},
	//	)
	//}
	//
	//// Set tokens to cookie
	//
	//cookie = new(http.Cookie)
	//cookie.Name = "access_token"
	//cookie.Value = *accessToken.Token
	//cookie.Path = "/"
	//cookie.MaxAge = int(a.Cfg.Jwt.AccessTokenMaxAge.Seconds())
	//c.SetCookie(cookie)
	//
	//return c.JSON(
	//	http.StatusOK,
	//	common_http.ResponseDTO{
	//		Status:  http.StatusOK,
	//		Message: "success",
	//		Data: map[string]interface{}{
	//			"access_token": accessToken.Token,
	//		},
	//	},
	//)

	return nil
}

// Login godoc
// @Summary Login user
// @Description Login user handler
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.Token
// @Router /auth/login [post]
func (a *AuthHandlers) Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var errorList []*common_http.ValidationError
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
				Error:   err.Error(),
			},
		)
	}

	accessToken, err := a.Jwt.CreateToken(meta.AccessToken.ID().String(), meta.UserID.String(), a.Cfg.Jwt.AccessTokenMaxAge)
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

	refreshToken, err := a.Jwt.CreateToken(meta.RefreshToken.ID().String(), meta.UserID.String(), a.Cfg.Jwt.RefreshTokenMaxAge)
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

	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = *accessToken.Token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = int(a.Cfg.Jwt.AccessTokenMaxAge.Seconds())
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = *refreshToken.Token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = int(a.Cfg.Jwt.RefreshTokenMaxAge.Seconds())
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
