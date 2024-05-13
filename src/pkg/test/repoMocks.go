package test

import (
	"context"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authz"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/events"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/bun/schema"
)

type MockBeforeAppendModelHook struct {
	mock.Mock
}

// BeforeAppendModel implements model.IBeforeAppendModelHook.
func (m *MockBeforeAppendModelHook) BeforeAppendModel(context.Context, schema.Query) error {
	panic("unimplemented")
}

var _ model.IBeforeAppendModelHook = &MockBeforeAppendModelHook{}

type MockGenericRepository[T any] struct {
	mock.Mock
}

func (mgr *MockGenericRepository[T]) Create(entity *T) error {
	args := mgr.Called(entity)

	if args.Get(0) != nil {
		args.Error(0)
	}

	return nil
}

func (mgr *MockGenericRepository[T]) GetById(id int) (*T, error) {
	args := mgr.Called(id)
	return args.Get(0).(*T), args.Error(1)
}

func (mgr *MockGenericRepository[T]) GetByUuid(uuid uuid.UUID) (*T, error) {
	args := mgr.Called(uuid)
	return args.Get(0).(*T), args.Error(1)
}

func (mgr *MockGenericRepository[T]) GetByField(field string, value interface{}) (*T, error) {
	args := mgr.Called(field, value)
	return args.Get(0).(*T), args.Error(1)
}

func (mgr *MockGenericRepository[T]) GetManyByField(field string, value interface{}) ([]T, error) {
	args := mgr.Called(field, value)
	return args.Get(0).([]T), args.Error(1)
}

func (mgr *MockGenericRepository[T]) Update(entity *T) error {
	args := mgr.Called(entity)
	return args.Error(0)
}

func (mgr *MockGenericRepository[T]) Delete(id int) error {
	args := mgr.Called(id)
	return args.Error(0)
}

func (mgr *MockGenericRepository[T]) ListAll() ([]T, error) {
	args := mgr.Called()
	return args.Get(0).([]T), args.Error(1)
}

type MockAuthzProvider struct {
	mock.Mock
}

// AddPermissionForRole implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) AddPermissionForRole(role authz.Role, permission authz.Permission) error {
	args := m.Called(role, permission)

	return args.Error(0)
}

// AddPermissionForRoleForResource implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) AddPermissionForRoleForResource(role authz.Role, resourceUuid uuid.UUID, permission authz.Permission) error {
	args := m.Called(role, resourceUuid, permission)

	return args.Error(0)
}

// AddPermissionForUser implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) AddPermissionForUser(userUuid uuid.UUID, permission authz.Permission) error {
	args := m.Called(userUuid, permission)

	return args.Error(0)
}

// AddRoleForUser implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) AddRoleForUser(userUuid uuid.UUID, role authz.Role) error {
	args := m.Called(userUuid, role)

	return args.Error(0)
}

// AddRoleForUserForResource implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) AddRoleForUserForResource(userUuid uuid.UUID, resourceUuid uuid.UUID, role authz.Role) error {
	args := m.Called(userUuid, resourceUuid, role)

	return args.Error(0)
}

// DeletePermissionForRole implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeletePermissionForRole(role authz.Role, permission authz.Permission) error {
	args := m.Called(role, permission)

	return args.Error(0)
}

// DeletePermissionForRoleForResource implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeletePermissionForRoleForResource(role authz.Role, resourceUuid uuid.UUID, permission authz.Permission) error {
	args := m.Called(role, resourceUuid, permission)

	return args.Error(0)
}

// DeletePermissionForUser implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeletePermissionForUser(userUuid uuid.UUID, permission authz.Permission) error {
	args := m.Called(userUuid, permission)

	return args.Error(0)
}

// DeleteRole implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeleteRole(role authz.Role) error {
	args := m.Called(role)

	return args.Error(0)
}

// DeleteRoleForResource implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeleteRoleForResource(role authz.Role, resourceUuid uuid.UUID) error {
	args := m.Called(role, resourceUuid)

	return args.Error(0)
}

// DeleteRoleForUser implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeleteRoleForUser(userUuid uuid.UUID, role authz.Role) error {
	args := m.Called(userUuid, role)

	return args.Error(0)
}

// DeleteRoleForUserForResource implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) DeleteRoleForUserForResource(userUuid uuid.UUID, resourceUuid uuid.UUID, role authz.Role) error {
	args := m.Called(userUuid, resourceUuid, role)

	return args.Error(0)
}

// GetMiddleware implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) GetMiddleware() echo.MiddlewareFunc {
	args := m.Called()

	return args.Get(0).(echo.MiddlewareFunc)
}

// GetPermissionConfig implements authz.AuthorizationProvider.
func (m *MockAuthzProvider) GetPermissionConfig() authz.PermissionConfig {
	args := m.Called()

	return args.Get(0).(authz.PermissionConfig)
}

var _ authz.AuthorizationProvider = &MockAuthzProvider{}

type MockEventManager struct {
	mock.Mock
}

// Emit implements events.EventManager.
func (m *MockEventManager) Emit(event events.Event) {
	m.Called(event)
}

// On implements events.EventManager.
// Subtle: this method shadows the method (Mock).On of MockEventManager.Mock.
func (m *MockEventManager) OnEvent(event events.Event, handler func(interface{})) {
	args := m.Called(event, handler)

	if args.Get(0) != nil {
		args.Get(0).(func())()
	}
}

var _ events.EventManager = &MockEventManager{}
