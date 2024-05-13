package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type IBeforeAppendModelHook interface {
	BeforeAppendModel(context.Context, bun.Query) error
}

type BaseModel struct {
	ID                    int                    `bun:"id,pk,autoincrement"`
	Uuid                  uuid.UUID              `bun:"uuid,type:uuid,default:gen_random_uuid()"`
	CreatedAt             time.Time              `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt             time.Time              `bun:",nullzero,notnull,default:current_timestamp"`
	BeforeAppendModelHook IBeforeAppendModelHook `bun:"-"`
}

func (b *BaseModel) UpdateTimestamps(query bun.Query) {
	switch query.(type) {
	case *bun.UpdateQuery:
		b.UpdatedAt = time.Now()
	}
}

type Account struct {
	BaseModel
	Username string
	About    string
}

func (a *Account) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	a.BaseModel.UpdateTimestamps(query)
	return nil
}

func NewAccount(username string, about string, userUuid uuid.UUID) *Account {
	return &Account{
		BaseModel: BaseModel{
			Uuid: userUuid,
		},
		Username: username,
		About:    about,
	}
}

type Notifcation struct {
	BaseModel
	AccountID int    `bun:"account_id,notnull"`
	Message   string `bun:"message,notnull"`
	IsRead    bool   `bun:"is_read,notnull"`
}

type Events struct {
	BaseModel
	ContractID         int      `bun:"contract_id,notnull"`
	ExecutingAccountID int      `bun:"executing_account_id,notnull"`
	ExecutingAccount   *Account `bun:"rel:belongs-to,join:executing_account_id=id"`
	Action             string   `bun:"action,notnull"`
}
