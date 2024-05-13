package repository

import "github.com/uptrace/bun"

type Repositories struct {
	AccountRepository IAccountRepository
}

func NewRepositories(bunDb *bun.DB) *Repositories {
	return &Repositories{
		AccountRepository: NewAccountRepository(bunDb),
	}
}
