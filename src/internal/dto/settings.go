package dto

import "github.com/VincentSchmid/htmx-templ-app-template/internal/model"

type UserEditData struct {
	Username string
	About    string
	Success  bool
}

func NewUserEditData(account *model.Account, success bool) UserEditData {
	return UserEditData{
		Username: account.Username,
		About:    account.About,
		Success:  success,
	}
}

type UserEditErrors struct {
	Username string
	About    string
}

type SecurityEditData struct {
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string
}

type SecurityEditErrors struct {
	ResetPassword string
	Success       bool
}

type FormData struct {
	UserSection     UserEditData
	SecuritySection SecurityEditData
}

type FormErrors struct {
	UserSection     UserEditErrors
	SecuritySection SecurityEditErrors
}

func NewFormErrors() FormErrors {
	return FormErrors{
		UserSection:     UserEditErrors{},
		SecuritySection: SecurityEditErrors{},
	}
}
