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

// Dataset types
const (
	DatasetQuiz          = "QUIZ"
	DatasetLoginApple    = "LOGIN_APPLE"
	DatasetSalesFLP      = "SALES_FLP"
	DatasetPointApple    = "POINT_APPLE"
	DatasetPointMyHero   = "POINT_MYHERO"
	DatasetTotalProspect = "TOTAL_PROSPECTS"
)

// Period frequencies
const (
	PeriodDaily     = "DAILY"
	PeriodWeekly    = "WEEKLY"
	PeriodMonthly   = "MONTHLY"
	PeriodQuarterly = "QUARTERLY"
	PeriodYearly    = "YEARLY"
)

var AllowedRoles = map[string]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleTL:         true,
	RoleViewer:     true,
	RoleSM:         true,
}

const (
	EntityTLActivity = "tl_activity"
	EntityTLCoaching = "tl_coaching"
	EntityTLBriefing = "tl_briefing"
)

const (
	SessionTypeCoaching = "coaching"
	SessionTypeBriefing = "briefing"
)
