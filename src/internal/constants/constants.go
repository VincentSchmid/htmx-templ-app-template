package constants

import "github.com/VincentSchmid/htmx-templ-app-template/pkg/authz"

type keyType string

const (
	USER_CONTEXT_KEY    keyType    = "user"
	UserRole            authz.Role = "user"
	AdminRole           authz.Role = "admin"
	ServiceProviderRole authz.Role = "serviceProvider"
	ClientRole          authz.Role = "client"
	ContractResource    string     = "contract"
)

const (
	GLOBAL_EVENT_CHANNEL = "global"

	CONTRACT_EVENT_CHANNEL = "contract"

	CONTRACT_UPDATED_EVENT = GLOBAL_EVENT_CHANNEL + ".contract.updated"
)
