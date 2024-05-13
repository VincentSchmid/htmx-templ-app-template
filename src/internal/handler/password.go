package handler

import (
	"log"

	authView "github.com/VincentSchmid/htmx-templ-app-template/internal/view/auth"
	uiView "github.com/VincentSchmid/htmx-templ-app-template/internal/view/ui"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authn"
	"github.com/labstack/echo/v4"
)

const (
	passwordResetPath = "/v1/auth/password/reset"
)

type PasswordHandler struct {
	authnProvider authn.AuthenticationProvider
}

var _ IHandler = (*PasswordHandler)(nil)

func NewPasswordHandler(authnProvider authn.AuthenticationProvider) *PasswordHandler {
	return &PasswordHandler{
		authnProvider: authnProvider,
	}
}

func (ph *PasswordHandler) sendResetRequest(c echo.Context) error {
	user := getAuthenticatedUser(c)
	redirectUrl := appconfig.Config.Webserver.Url + passwordResetPath
	err := ph.authnProvider.ResetPassword(c, user.Email, redirectUrl)
	if err != nil {
		log.Printf("failed to reset password: %v", err)
		return render(c, uiView.ErrorLabel(err.Error()))
	}

	return render(c, authView.ForgotPasswordInitiated(user.Email))
}

func (ph *PasswordHandler) passwordResetCallback(c echo.Context) error {
	return render(c, authView.ResetPasswordIndex())
}

func (ph *PasswordHandler) setNewPassword(c echo.Context) error {
	credentials := authn.UserCredentials{
		Password:        c.FormValue("password"),
		ConfirmPassword: c.FormValue("confirmPassword"),
	}

	user := getAuthenticatedUser(c)

	err := ph.authnProvider.ChangePassword(c, user.AccessToken, credentials)
	if err != nil {
		return render(c, authView.ResetPasswordForm("password reset failed"))
	}

	return hxRedirect(c, "/")
}

func (ph *PasswordHandler) forgotPasswordIndex(c echo.Context) error {
	return render(c, authView.ForgotPasswordIndex())
}

func (ph *PasswordHandler) forgotPassword(c echo.Context) error {
	email := c.FormValue("email")

	err := ph.authnProvider.ResetPassword(c, email, appconfig.Config.Webserver.Url+passwordResetPath)
	if err != nil {
		log.Printf("failed to reset password: %v", err)
		return render(c, uiView.ErrorLabel(err.Error()))
	}

	return render(c, authView.ForgotPasswordInitiated(email))
}

func (ph *PasswordHandler) RegisterRoutes(basePath string, e *echo.Group) {
	group := e.Group(basePath)
	group.GET("/forgot", ph.forgotPasswordIndex)
	group.POST("/forgot", ph.forgotPassword)
	group.GET("/callback", ph.passwordResetCallback)
	group.GET("/reset", ph.setNewPassword)
	group.POST("/reset", ph.sendResetRequest)
}
