package service

import (
	"log"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/constants"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/dto"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/repository"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authz"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AccountService struct {
	accountRepository repository.IAccountRepository
	authZProvider     authz.AuthorizationProvider
}

func NewAccountService(accountRepository repository.IAccountRepository, authZProvider authz.AuthorizationProvider) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		authZProvider:     authZProvider,
	}
}

func (as *AccountService) GetAccount(userUuid uuid.UUID) model.Account {
	account, err := as.accountRepository.BasicOperations().GetByUuid(userUuid)
	if err != nil {
		return model.Account{}
	}

	return *account
}

func (as *AccountService) CreateOrUpdateAccount(form dto.UserEditData, userUuid uuid.UUID) (model.Account, error) {
	account := model.NewAccount(form.Username, form.About, userUuid)

	if err := as.accountRepository.CreateOrUpdate(account); err != nil {
		logger.Log.Error("failed to update or create account: ",
			zap.Error(err))

		return model.Account{}, ErrInternal
	}

	err := as.authZProvider.AddRoleForUser(account.Uuid, constants.UserRole)
	log.Println("Added role for user")
	if err != nil {
		logger.Log.Error("failed to add role for user: ", zap.Error(err))
		return model.Account{}, ErrInternal
	}

	return *account, nil
}
