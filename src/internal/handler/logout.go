package handler

import (
	"fmt"

	"github.com/VincentSchmid/htmx-templ-app-template/pkg/session"
	"github.com/labstack/echo/v4"
)

type LogoutHandler struct {
	// Any dependencies, such as a service for user authentication, can be included here.
}

var _ IHandler = (*LogoutHandler)(nil)

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

func (lh *LogoutHandler) logoutCreate(c echo.Context) error {
	err := session.ClearAccessToken(c)
	if err != nil {
		return fmt.Errorf("failed to clear access token: %w", err)
	}

	return hxRedirect(c, "/")
}

func (lh *LogoutHandler) RegisterRoutes(basePath string, e *echo.Group) {
	e.POST(basePath, lh.logoutCreate)
}
