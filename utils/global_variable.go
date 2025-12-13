package utils

const (
	CtxKeyId       = "CTX_ID"
	CtxKeyAuthData = "auth_data"
)

const (
	RedisAppConf = "cache:config:app"
	RedisDbConf  = "cache:config:db"
)

const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleTL         = "teamleader"
	RoleSM         = "salesman"
	RoleViewer     = "viewer"
)

var AllowedRoles = map[string]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleTL:         true,
	RoleViewer:     true,
	RoleSM:         true,
}
