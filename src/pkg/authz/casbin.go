package authz

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	casbin_mw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CasbinAuthzProvider struct {
	Enforcer         *casbin.Enforcer
	UserGetter       func(c echo.Context) (string, error)
	permissionConfig PermissionConfig
	policies         []*appconfig.Policy
	permissionMap    map[*regexp.Regexp]string
}

var _ AuthorizationProvider = (*CasbinAuthzProvider)(nil)

func NewCasbinAuthzProvider(config *appconfig.Casbin, userGetter func(c echo.Context) (string, error), db *sql.DB) *CasbinAuthzProvider {
	err := prepopulatePolicies(db, config.DefaultPolicies, config.TableName)
	if err != nil {
		panic(err)
	}

	adapter, err := sqladapter.NewAdapter(db, config.DbDriver, config.TableName)
	if err != nil {
		panic(err)
	}

	model, err := model.NewModelFromString(config.Model)
	if err != nil {
		panic(err)
	}

	var permissionMap = make(map[*regexp.Regexp]string)
	for pattern, permission := range config.ResourceByPathRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		permissionMap[re] = permission
	}

	enforcer, err := casbin.NewEnforcer(model, adapter)
	if err != nil {
		panic(err)
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}

	return &CasbinAuthzProvider{
		Enforcer:      enforcer,
		UserGetter:    userGetter,
		permissionMap: permissionMap,
		policies:      config.Policies,
		permissionConfig: PermissionConfig{
			Read:      config.PermissionNames.Reader,
			Write:     config.PermissionNames.Writer,
			ReadWrite: config.PermissionNames.ReaderWriter,
		},
	}
}

func (crp *CasbinAuthzProvider) GetPermissionConfig() PermissionConfig {
	return crp.permissionConfig
}

func (crp *CasbinAuthzProvider) AddRoleForUser(userUuid uuid.UUID, role Role) error {
	_, err := crp.Enforcer.AddRoleForUser(userUuid.String(), string(role))
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) DeleteRoleForUser(userUuid uuid.UUID, role Role) error {
	_, err := crp.Enforcer.DeleteRoleForUser(userUuid.String(), string(role))
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) DeleteRole(role Role) error {
	_, err := crp.Enforcer.DeleteRole(string(role))
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) AddPermissionForRole(role Role, permission Permission) error {
	_, err := crp.Enforcer.AddPolicy(string(role), permission.Resource, permission.Method)
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) DeletePermissionForRole(role Role, permission Permission) error {
	_, err := crp.Enforcer.RemovePolicy(string(role), permission.Resource, permission.Method)
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) AddPermissionForUser(userUuid uuid.UUID, permission Permission) error {
	_, err := crp.Enforcer.AddPolicy(userUuid.String(), permission.Resource, permission.Method)
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) DeletePermissionForUser(userUuid uuid.UUID, permission Permission) error {
	_, err := crp.Enforcer.RemovePolicy(userUuid.String(), permission.Resource, permission.Method)
	if err != nil {
		return err
	}
	return crp.Enforcer.SavePolicy()
}

func (crp *CasbinAuthzProvider) AddRoleForUserForResource(userUuid uuid.UUID, resourceUuid uuid.UUID, role Role) error {
	logger.Log.Info("adding role for user for resource", zap.String("user_uuid", userUuid.String()), zap.String("resource_uuid", resourceUuid.String()), zap.String("role", string(role)))
	resourceScopedRole := getResourceRole(role, resourceUuid)
	err := crp.addResourceScopedPermissions(role, resourceUuid)
	if err != nil {
		return fmt.Errorf("failed to add resource scoped permissions: %w", err)
	}

	return crp.AddRoleForUser(userUuid, resourceScopedRole)
}

func (crp *CasbinAuthzProvider) DeleteRoleForUserForResource(userUuid uuid.UUID, resourceUuid uuid.UUID, role Role) error {
	resourceScopedRole := getResourceRole(role, resourceUuid)
	return crp.DeleteRoleForUser(userUuid, resourceScopedRole)
}

func (crp *CasbinAuthzProvider) DeleteRoleForResource(role Role, resourceUuid uuid.UUID) error {
	resourceScopedRole := getResourceRole(role, resourceUuid)
	return crp.DeleteRole(resourceScopedRole)
}

func (crp *CasbinAuthzProvider) AddPermissionForRoleForResource(role Role, resourceUuid uuid.UUID, permission Permission) error {
	resourceScopedRole := getResourceRole(role, resourceUuid)
	permission.Resource = getResourcePermission(permission.Resource, resourceUuid)
	return crp.AddPermissionForRole(resourceScopedRole, permission)
}

func (crp *CasbinAuthzProvider) DeletePermissionForRoleForResource(role Role, resourceUuid uuid.UUID, permission Permission) error {
	resourceScopedRole := getResourceRole(role, resourceUuid)
	permission.Resource = getResourcePermission(permission.Resource, resourceUuid)
	return crp.DeletePermissionForRole(resourceScopedRole, permission)
}

func (crp *CasbinAuthzProvider) GetMiddleware() echo.MiddlewareFunc {
	return casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		Enforcer:       crp.Enforcer,
		UserGetter:     crp.UserGetter,
		EnforceHandler: crp.enforceHandler,
	})
}

func (crp *CasbinAuthzProvider) enforceHandler(c echo.Context, user string) (bool, error) {
	method := c.Request().Method
	path := c.Request().URL.Path
	permission, ok := crp.findPermissionForPath(path)
	if !ok {
		return false, fmt.Errorf("forbidden")
	}

	ok, err := crp.Enforcer.Enforce(user, permission, method)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("forbidden")
	}

	return true, nil
}

func (crp *CasbinAuthzProvider) findPermissionForPath(path string) (string, bool) {
	for re, permission := range crp.permissionMap {
		if matches := re.FindStringSubmatch(path); matches != nil {
			if strings.Contains(permission, "{uuid}") && len(matches) > 1 {
				return strings.Replace(permission, "{uuid}", matches[1], 1), true
			}
			return permission, true
		}
	}

	return "", false
}

func (crp *CasbinAuthzProvider) addResourceScopedPermissions(role Role, resourceUuid uuid.UUID) error {
	resourceScopedRole := getResourceRole(role, resourceUuid)
	roleString := string(role) + ":{uuid}"
	for _, policy := range crp.policies {
		if policy.Role == roleString {
			for resource, permission := range policy.Permissions {
				resourceScopedPermission := Permission{
					Resource: getResourcePermission(resource, resourceUuid),
					Method:   permission,
				}
				err := crp.AddPermissionForRole(resourceScopedRole, resourceScopedPermission)
				if err != nil {
					return fmt.Errorf("failed to add permission for role: %w", err)
				}
			}
		}
	}

	return nil
}

func getResourceRole(role Role, resourceUuid uuid.UUID) Role {
	return Role(string(role) + ":" + resourceUuid.String())
}

func getResourcePermission(resource string, resourceUuid uuid.UUID) string {
	return strings.Replace(resource, "{uuid}", resourceUuid.String(), 1)
}

func prepopulatePolicies(db *sql.DB, policies []*appconfig.Policy, tableName string) error {
	for _, policy := range policies {
		for action, subject := range policy.Permissions {

			var exists bool
			queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE p_type = $1 AND v0 = $2 AND v1 = $3 AND v2 = $4 LIMIT 1)", tableName)
			err := db.QueryRow(queryCheck, "p", policy.Role, action, subject).Scan(&exists)
			if err != nil {
				return err
			}

			if !exists {
				queryInsert := fmt.Sprintf("INSERT INTO %s (p_type, v0, v1, v2) VALUES ($1, $2, $3, $4)", tableName)
				_, err := db.Exec(queryInsert, "p", policy.Role, action, subject)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
