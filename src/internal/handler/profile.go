package handler

import (
	"github.com/VincentSchmid/htmx-templ-app-template/internal/dto"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/service"
	view "github.com/VincentSchmid/htmx-templ-app-template/internal/view/settings"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	accountService *service.AccountService
}

var _ IHandler = (*UserHandler)(nil)

func NewUserHandler(accountService *service.AccountService) *UserHandler {
	return &UserHandler{
		accountService: accountService,
	}
}

func (uh *UserHandler) accountSetupIndex(c echo.Context) error {
	user := getAuthenticatedUser(c)
	account := uh.accountService.GetAccount(user.Uuid)

	formData := dto.FormData{
		UserSection:     dto.NewUserEditData(&account, false),
		SecuritySection: dto.SecurityEditData{},
	}

	formErrors := dto.NewFormErrors()

	return render(c, view.ProfileSetupIndex(formData, formErrors))
}

func (uh *UserHandler) updateProfile(c echo.Context) error {
	shouldRedirect := true
	redirect, err := GetRedirectUrl(c)
	if err != nil {
		shouldRedirect = false
	}

	params := dto.UserEditData{
		Username: c.FormValue("username"),
		About:    c.FormValue("about"),
	}

	user := getAuthenticatedUser(c)

	account, err := uh.accountService.CreateOrUpdateAccount(params, user.Uuid)
	if err != nil {
		return render(c, view.UserEditForm(dto.NewUserEditData(&account, false), dto.UserEditErrors{}))
	}

	if shouldRedirect {
		return hxRedirect(c, redirect)
	}

	return render(c, view.UserEditForm(dto.NewUserEditData(&account, true), dto.UserEditErrors{}))
}

func (uh *UserHandler) settingsIndex(c echo.Context) error {
	user := getAuthenticatedUser(c)
	return render(c, view.Index(user))
}

func (uh *UserHandler) RegisterRoutes(basePath string, e *echo.Group) {
	group := e.Group(basePath)
	group.GET("/profile", uh.accountSetupIndex)
	group.PUT("/profile", uh.updateProfile)
	group.GET("/settings", uh.settingsIndex)
}
