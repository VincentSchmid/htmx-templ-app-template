package model

import (
	"github.com/google/uuid"
)

type AuthenticatedUser struct {
	Uuid        uuid.UUID
	Email       string
	IsLoggedIn  bool
	Account     *Account
	AccessToken string
}

func NewAuthenticatedUser() *AuthenticatedUser {
	return &AuthenticatedUser{Account: &Account{}}
}
