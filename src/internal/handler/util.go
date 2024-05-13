package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/constants"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	uiView "github.com/VincentSchmid/htmx-templ-app-template/internal/view/ui"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func hxRedirect(c echo.Context, url string) error {
	if len(c.Request().Header.Get("Hx-Request")) > 0 {
		c.Response().Header().Set("Hx-Redirect", url)
		c.Response().WriteHeader(http.StatusSeeOther)
		return nil
	}

	return c.Redirect(http.StatusSeeOther, url)
}

func renderError(c echo.Context, message string) error {
	c.Response().WriteHeader(http.StatusInternalServerError)

	notificationData := uiView.NotificationData{
		Type:    uiView.NotificationTypeError,
		Title:   "Internal Server Error",
		Details: message,
	}

	return render(c, uiView.Notification(notificationData))
}

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func Make(h func(c echo.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := h(c); err != nil {
			zap.L().Error("internal server error", zap.Error(err))
			return err
		}
		return nil
	}
}

func getAuthenticatedUser(c echo.Context) model.AuthenticatedUser {
	user, ok := c.Request().Context().Value(constants.USER_CONTEXT_KEY).(model.AuthenticatedUser)
	if !ok {
		return *model.NewAuthenticatedUser()
	}
	return user
}

func GetUserIdFunc() func(c echo.Context) (string, error) {
	return func(c echo.Context) (string, error) {
		user, ok := c.Request().Context().Value(constants.USER_CONTEXT_KEY).(model.AuthenticatedUser)
		if !ok {
			return "", fmt.Errorf("failed to get authenticated user")
		}

		return user.Uuid.String(), nil
	}
}

func GetRedirectUrl(c echo.Context) (string, error) {
	referer := c.Request().Header.Get("Referer")

	parsedURL, err := url.Parse(referer)
	if err != nil {
		return "", fmt.Errorf("failed to parse referer URL: %w", err)
	}

	queries := parsedURL.Query()

	redirect := queries.Get("redirect")
	if redirect == "" {
		return "", fmt.Errorf("failed to get redirect URL from referer")
	}

	return redirect, nil
}
