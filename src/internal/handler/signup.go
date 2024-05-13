package handler

import (
	"fmt"

	authView "github.com/VincentSchmid/htmx-templ-app-template/internal/view/auth"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authn"
	"github.com/labstack/echo/v4"
)

type SignUpHandler struct {
	authnProvider authn.AuthenticationProvider
}

var _ IHandler = (*SignUpHandler)(nil)

func NewSignUpHandler(authnProvider authn.AuthenticationProvider) *SignUpHandler {
	return &SignUpHandler{
		authnProvider: authnProvider,
	}
}

func (sh *SignUpHandler) signUpIndex(c echo.Context) error {
	return render(c, authView.SignUpIndex())
}

func (sh *SignUpHandler) signUpCreate(c echo.Context) error {
	userCredentials := authn.UserCredentials{
		Email:           c.FormValue("email"),
		Password:        c.FormValue("password"),
		ConfirmPassword: c.FormValue("confirmPassword"),
	}

	user, err := sh.authnProvider.SignUpWithCredentials(c, userCredentials)
	if err != nil {
		return fmt.Errorf("failed to sign up: %w", err)
	}

	return render(c, authView.SignupSuccess(user.Email))
}

func (sh *SignUpHandler) RegisterRoutes(basePath string, e *echo.Group) {
	e.GET(basePath, sh.signUpIndex)
	e.POST(basePath, sh.signUpCreate)
}
