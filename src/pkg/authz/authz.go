package authz

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Role string

type Permission struct {
	Resource string
	Method   string
}

type AuthorizationProvider interface {
	GetPermissionConfig() PermissionConfig
	AddRoleForUser(userUuid uuid.UUID, role Role) error
	DeleteRoleForUser(userUuid uuid.UUID, role Role) error
	DeleteRole(role Role) error
	AddPermissionForRole(role Role, permission Permission) error
	DeletePermissionForRole(role Role, permission Permission) error
	AddPermissionForUser(userUuid uuid.UUID, permission Permission) error
	DeletePermissionForUser(userUuid uuid.UUID, permission Permission) error

	AddRoleForUserForResource(userUuid uuid.UUID, resourceUuid uuid.UUID, role Role) error
	DeleteRoleForUserForResource(userUuid uuid.UUID, resourceUuid uuid.UUID, role Role) error
	DeleteRoleForResource(role Role, resourceUuid uuid.UUID) error
	AddPermissionForRoleForResource(role Role, resourceUuid uuid.UUID, permission Permission) error
	DeletePermissionForRoleForResource(role Role, resourceUuid uuid.UUID, permission Permission) error

	GetMiddleware() echo.MiddlewareFunc
}

type PermissionConfig struct {
	Read      string
	Write     string
	ReadWrite string
}
