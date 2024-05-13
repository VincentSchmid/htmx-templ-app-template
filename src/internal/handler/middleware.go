package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/constants"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/repository"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authn"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/session"
	"github.com/labstack/echo/v4"
)

const (
	loginPath   = "/v1/auth/login"
	profilePath = "/v1/user/profile"
)

func GetWithUserMiddleware(authnProvider authn.AuthenticationProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/public/") {
				return next(c)
			}
			accessToken, err := session.GetAccessToken(c)
			if err != nil {
				return next(c)
			}

			user, err := authnProvider.AuthenticateWithToken(c, accessToken)
			if err != nil {
				log.Println("failed to authenticate with token:", err)
				return next(c)
			}

			ctx := context.WithValue(c.Request().Context(), constants.USER_CONTEXT_KEY, user)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func WithAuthn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/public/") {
			return next(c)
		}
		user := getAuthenticatedUser(c)
		if !user.IsLoggedIn {
			return c.Redirect(http.StatusSeeOther, loginPath+"?redirect="+c.Request().URL.Path)
		}

		return next(c)
	}
}
func GetWithAccountSetupMiddleware(accountRepository repository.IAccountRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := getAuthenticatedUser(c)
			account, err := accountRepository.BasicOperations().GetByUuid(user.Uuid)
			if errors.Is(err, sql.ErrNoRows) {
				return c.Redirect(http.StatusSeeOther, profilePath+"?redirect="+c.Request().URL.Path)
			}

			if err != nil {
				return fmt.Errorf("failed to get account: %w", err)
			}

			user.Account = account
			ctx := context.WithValue(c.Request().Context(), constants.USER_CONTEXT_KEY, user)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
