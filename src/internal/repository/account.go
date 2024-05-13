package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	"github.com/uptrace/bun"
)

type IAccountRepository interface {
	BasicOperations() IGenericRepository[model.Account]
	CreateOrUpdate(account *model.Account) error
}

type AccountRepository struct {
	bun             *bun.DB
	ctx             context.Context
	basicOperations IGenericRepository[model.Account]
}

var _ IAccountRepository = (*AccountRepository)(nil)

func NewAccountRepository(bunDb *bun.DB) *AccountRepository {
	return &AccountRepository{
		bun:             bunDb,
		ctx:             context.Background(),
		basicOperations: NewGenericRepository[model.Account](bunDb),
	}
}

func (ar *AccountRepository) CreateOrUpdate(account *model.Account) error {
	existingAccount, err := ar.basicOperations.GetByUuid(account.Uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return ar.basicOperations.Create(account)
	}
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	account.BaseModel = existingAccount.BaseModel
	err = ar.basicOperations.Update(account)

	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func (ar *AccountRepository) BasicOperations() IGenericRepository[model.Account] {
	return ar.basicOperations
}
