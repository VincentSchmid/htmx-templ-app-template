package handler

import (
	"github.com/VincentSchmid/htmx-templ-app-template/internal/service"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authn"
)

type Handlers struct {
	LoginHandler    *LoginHandler
	LogoutHandler   *LogoutHandler
	PasswordHandler *PasswordHandler
	UserHandler     *UserHandler
	SignUpHandler   *SignUpHandler
}

func NewHandlers(services *service.Services, authnProvider authn.AuthenticationProvider) *Handlers {
	return &Handlers{
		LoginHandler:    NewLoginHandler(authnProvider),
		LogoutHandler:   NewLogoutHandler(),
		PasswordHandler: NewPasswordHandler(authnProvider),
		UserHandler:     NewUserHandler(services.AccountService),
		SignUpHandler:   NewSignUpHandler(authnProvider),
	}
}
