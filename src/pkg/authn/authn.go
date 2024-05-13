package authn

import (
	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserCredentials struct {
	Email           string
	Password        string
	ConfirmPassword string
}

type User struct {
	Uuid  uuid.UUID
	Email string
}

type OAuthUrl string

type AuthenticationProvider interface {
	SignUpWithCredentials(context echo.Context, credentials UserCredentials) (User, error)
	SignInWithCredentials(context echo.Context, credentials UserCredentials) (string, error)
	SignInWithProvider(provider string, redirectUrl string) (OAuthUrl, error)
	AuthenticateWithToken(context echo.Context, token string) (model.AuthenticatedUser, error)
	ChangePassword(context echo.Context, token string, newCredentials UserCredentials) error
	ResetPassword(context echo.Context, email string, redirectUrl string) error
}
