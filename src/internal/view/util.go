package view

import (
	"context"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/constants"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
)

func GetAuthenticatedUser(ctx context.Context) model.AuthenticatedUser {
	user, ok := ctx.Value(constants.USER_CONTEXT_KEY).(model.AuthenticatedUser)
	if !ok {
		return *model.NewAuthenticatedUser()
	}
	return user
}
