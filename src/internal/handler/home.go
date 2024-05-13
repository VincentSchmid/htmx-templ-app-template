package handler

import (
	homeView "github.com/VincentSchmid/htmx-templ-app-template/internal/view/home"
	"github.com/labstack/echo/v4"
)

func HandleHomeIndex(c echo.Context) error {
	user := getAuthenticatedUser(c)

	if !user.IsLoggedIn {
		return render(c, homeView.Index())
	}

	return render(c, homeView.Index())
}
