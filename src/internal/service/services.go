package service

import (
	"github.com/VincentSchmid/htmx-templ-app-template/internal/repository"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authz"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/events"
)

type Services struct {
	AccountService *AccountService
}

func NewServices(respositories *repository.Repositories,
	authZProvider authz.AuthorizationProvider,
	eventManager events.EventManager) *Services {
	return &Services{
		AccountService: NewAccountService(respositories.AccountRepository, authZProvider),
	}
}
