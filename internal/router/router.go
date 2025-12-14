package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"teamleader-management/infrastructure/database"
	datasetHandler "teamleader-management/internal/handlers/http/dataset"
	kpiHandler "teamleader-management/internal/handlers/http/kpiitem"
	menuHandler "teamleader-management/internal/handlers/http/menu"
	permissionHandler "teamleader-management/internal/handlers/http/permission"
	personHandler "teamleader-management/internal/handlers/http/person"
	pillarHandler "teamleader-management/internal/handlers/http/pillar"
	roleHandler "teamleader-management/internal/handlers/http/role"
	sessionHandler "teamleader-management/internal/handlers/http/session"
	tlHandler "teamleader-management/internal/handlers/http/tl"
	userHandler "teamleader-management/internal/handlers/http/user"
	authRepo "teamleader-management/internal/repositories/auth"
	datasetRepo "teamleader-management/internal/repositories/dataset"
	kpiRepo "teamleader-management/internal/repositories/kpiitem"
	mediaRepo "teamleader-management/internal/repositories/media"
	menuRepo "teamleader-management/internal/repositories/menu"
	metricRepo "teamleader-management/internal/repositories/metric"
	permissionRepo "teamleader-management/internal/repositories/permission"
	personRepo "teamleader-management/internal/repositories/person"
	pillarRepo "teamleader-management/internal/repositories/pillar"
	roleRepo "teamleader-management/internal/repositories/role"
	sessionRepo "teamleader-management/internal/repositories/session"
	tlActivityRepo "teamleader-management/internal/repositories/tlactivity"
	tlAttendanceRepo "teamleader-management/internal/repositories/tlattendance"
	tlSessionRepo "teamleader-management/internal/repositories/tlsession"
	tlTrainingRepo "teamleader-management/internal/repositories/tltraining"
	userRepo "teamleader-management/internal/repositories/user"
	datasetSvc "teamleader-management/internal/services/dataset"
	kpiSvc "teamleader-management/internal/services/kpiitem"
	mediaSvc "teamleader-management/internal/services/media"
	menuSvc "teamleader-management/internal/services/menu"
	permissionSvc "teamleader-management/internal/services/permission"
	personSvc "teamleader-management/internal/services/person"
	pillarSvc "teamleader-management/internal/services/pillar"
	roleSvc "teamleader-management/internal/services/role"
	sessionSvc "teamleader-management/internal/services/session"
	tlActivitySvc "teamleader-management/internal/services/tlactivity"
	tlAttendanceSvc "teamleader-management/internal/services/tlattendance"
	tlSessionSvc "teamleader-management/internal/services/tlsession"
	tlTrainingSvc "teamleader-management/internal/services/tltraining"
	userSvc "teamleader-management/internal/services/user"
	"teamleader-management/middlewares"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/security"
	"teamleader-management/utils"
)

type Routes struct {
	App *gin.Engine
	DB  *gorm.DB
}

func NewRoutes() *Routes {
	app := gin.Default()

	app.Use(middlewares.CORS())
	app.Use(gin.CustomRecovery(middlewares.ErrorHandler))
	app.Use(middlewares.SetContextId())

	app.GET("/healthcheck", func(ctx *gin.Context) {
		logger.WriteLog(logger.LogLevelDebug, "ClientIP: "+ctx.ClientIP())
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK!!",
		})
	})

	return &Routes{
		App: app,
	}
}

func (r *Routes) UserRoutes() {
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	repo := userRepo.NewUserRepo(r.DB)
	rRepo := roleRepo.NewRoleRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	prRepo := personRepo.NewPersonRepo(r.DB)
	uc := userSvc.NewUserService(repo, blacklistRepo, rRepo, pRepo, prRepo)

	// Setup login limiter if Redis is available
	redisClient := database.GetRedisClient()
	var loginLimiter security.LoginLimiter
	if redisClient != nil {
		loginLimiter = security.NewRedisLoginLimiter(
			redisClient,
			utils.GetEnv("LOGIN_ATTEMPT_LIMIT", 5).(int),
			time.Duration(utils.GetEnv("LOGIN_ATTEMPT_WINDOW_SECONDS", 60).(int))*time.Second,
			time.Duration(utils.GetEnv("LOGIN_BLOCK_DURATION_SECONDS", 300).(int))*time.Second,
		)
	}

	h := userHandler.NewUserHandler(uc, loginLimiter)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// Setup register rate limiter
	registerLimit := utils.GetEnv("REGISTER_RATE_LIMIT", 5).(int)
	registerWindowSeconds := utils.GetEnv("REGISTER_RATE_WINDOW_SECONDS", 60).(int)
	if registerWindowSeconds <= 0 {
		registerWindowSeconds = 60
	}
	registerLimiter := middlewares.IPRateLimitMiddleware(
		redisClient,
		"user_register",
		registerLimit,
		time.Duration(registerWindowSeconds)*time.Second,
	)

	user := r.App.Group("/api/user")
	{
		user.POST("/register", registerLimiter, h.Register)
		user.POST("/login", h.Login)
		user.POST("/forgot-password", h.ForgotPassword)
		user.POST("/reset-password", h.ResetPassword)

		userPriv := user.Group("").Use(mdw.AuthMiddleware())
		{
			userPriv.POST("/logout", h.Logout)
			userPriv.GET("", h.GetUserByAuth)
			userPriv.GET("/:id", mdw.PermissionMiddleware("users", "view"), h.GetUserById)
			userPriv.PUT("", h.Update)
			userPriv.PUT("/:id", mdw.PermissionMiddleware("users", "update"), h.UpdateUserById)
			userPriv.PUT("/change/password", h.ChangePassword)
			userPriv.DELETE("", h.Delete)
			userPriv.DELETE("/:id", mdw.PermissionMiddleware("users", "delete"), h.DeleteUserById)

			// Admin create user endpoint (with role selection)
			userPriv.POST("", mdw.PermissionMiddleware("users", "create"), h.AdminCreateUser)
		}
	}

	r.App.GET("/api/users", mdw.AuthMiddleware(), mdw.PermissionMiddleware("users", "list"), h.GetAllUsers)
}

func (r *Routes) RoleRoutes() {
	repoRole := roleRepo.NewRoleRepo(r.DB)
	repoPermission := permissionRepo.NewPermissionRepo(r.DB)
	repoMenu := menuRepo.NewMenuRepo(r.DB)
	svc := roleSvc.NewRoleService(repoRole, repoPermission, repoMenu)
	h := roleHandler.NewRoleHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, repoPermission)

	// List endpoints
	r.App.GET("/api/roles", mdw.AuthMiddleware(), mdw.PermissionMiddleware("roles", "list"), h.GetAll)

	// CRUD endpoints
	role := r.App.Group("/api/role").Use(mdw.AuthMiddleware())
	{
		role.POST("", mdw.PermissionMiddleware("roles", "create"), h.Create)
		role.GET("/:id", mdw.PermissionMiddleware("roles", "view"), h.GetByID)
		role.PUT("/:id", mdw.PermissionMiddleware("roles", "update"), h.Update)
		role.DELETE("/:id", mdw.PermissionMiddleware("roles", "delete"), h.Delete)

		// Permission and menu assignment
		role.POST("/:id/permissions", mdw.PermissionMiddleware("roles", "assign_permissions"), h.AssignPermissions)
		role.POST("/:id/menus", mdw.PermissionMiddleware("roles", "assign_menus"), h.AssignMenus)
	}
}

func (r *Routes) PermissionRoutes() {
	repo := permissionRepo.NewPermissionRepo(r.DB)
	svc := permissionSvc.NewPermissionService(repo)
	h := permissionHandler.NewPermissionHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, repo)

	// List endpoints
	r.App.GET("/api/permissions", mdw.AuthMiddleware(), mdw.PermissionMiddleware("permissions", "list"), h.GetAll)

	// Get current user's permissions
	r.App.GET("/api/permissions/me", mdw.AuthMiddleware(), h.GetUserPermissions)

	// CRUD endpoints
	permission := r.App.Group("/api/permission").Use(mdw.AuthMiddleware())
	{
		permission.POST("", mdw.PermissionMiddleware("permissions", "create"), h.Create)
		permission.GET("/:id", mdw.PermissionMiddleware("permissions", "view"), h.GetByID)
		permission.PUT("/:id", mdw.PermissionMiddleware("permissions", "update"), h.Update)
		permission.DELETE("/:id", mdw.PermissionMiddleware("permissions", "delete"), h.Delete)
	}
}

func (r *Routes) MenuRoutes() {
	repo := menuRepo.NewMenuRepo(r.DB)
	svc := menuSvc.NewMenuService(repo)
	h := menuHandler.NewMenuHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// Public endpoints for authenticated users
	r.App.GET("/api/menus/active", mdw.AuthMiddleware(), h.GetActiveMenus)
	r.App.GET("/api/menus/me", mdw.AuthMiddleware(), h.GetUserMenus)

	// List endpoints
	r.App.GET("/api/menus", mdw.AuthMiddleware(), mdw.PermissionMiddleware("menus", "list"), h.GetAll)

	// CRUD endpoints
	menu := r.App.Group("/api/menu").Use(mdw.AuthMiddleware())
	{
		menu.POST("", mdw.PermissionMiddleware("menus", "create"), h.Create)
		menu.GET("/:id", mdw.PermissionMiddleware("menus", "view"), h.GetByID)
		menu.PUT("/:id", mdw.PermissionMiddleware("menus", "update"), h.Update)
		menu.DELETE("/:id", mdw.PermissionMiddleware("menus", "delete"), h.Delete)
	}
}

func (r *Routes) PillarRoutes() {
	repo := pillarRepo.NewPillarRepo(r.DB)
	svc := pillarSvc.NewPillarService(repo)
	h := pillarHandler.NewPillarHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	r.App.GET("/api/pillars", mdw.AuthMiddleware(), mdw.PermissionMiddleware("pillars", "list"), h.GetAll)

	pillar := r.App.Group("/api/pillar").Use(mdw.AuthMiddleware())
	{
		pillar.POST("", mdw.PermissionMiddleware("pillars", "create"), h.Create)
		pillar.GET("/:id", mdw.PermissionMiddleware("pillars", "view"), h.GetByID)
		pillar.PUT("/:id", mdw.PermissionMiddleware("pillars", "update"), h.Update)
		pillar.DELETE("/:id", mdw.PermissionMiddleware("pillars", "delete"), h.Delete)
	}
}

func (r *Routes) KPIItemRoutes() {
	kRepo := kpiRepo.NewKPIItemRepo(r.DB)
	tRepo := kpiRepo.NewPersonKPITargetRepo(r.DB)
	pRepo := pillarRepo.NewPillarRepo(r.DB)
	prRepo := personRepo.NewPersonRepo(r.DB)
	svc := kpiSvc.NewKPIItemService(kRepo, pRepo, prRepo, tRepo)
	h := kpiHandler.NewKPIItemHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	permRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, permRepo)

	r.App.GET("/api/kpi-items", mdw.AuthMiddleware(), mdw.PermissionMiddleware("kpi_items", "list"), h.GetAll)

	kpi := r.App.Group("/api/kpi-item").Use(mdw.AuthMiddleware())
	{
		kpi.POST("", mdw.PermissionMiddleware("kpi_items", "create"), h.Create)
		kpi.GET("/:id", mdw.PermissionMiddleware("kpi_items", "view"), h.GetByID)
		kpi.PUT("/:id", mdw.PermissionMiddleware("kpi_items", "update"), h.Update)
		kpi.DELETE("/:id", mdw.PermissionMiddleware("kpi_items", "delete"), h.Delete)
		kpi.POST("/target", mdw.PermissionMiddleware("kpi_items", "update"), h.UpsertPersonTarget)
		kpi.DELETE("/target", mdw.PermissionMiddleware("kpi_items", "delete"), h.DeletePersonTarget)
	}
}

func (r *Routes) DatasetRoutes() {
	repo := datasetRepo.NewDatasetRepo(r.DB)
	mRepo := metricRepo.NewMetricRepo(r.DB)
	pRepo := personRepo.NewPersonRepo(r.DB)
	svc := datasetSvc.NewDatasetService(repo)
	processor := datasetSvc.NewProcessor(repo, mRepo, pRepo)
	h := datasetHandler.NewDatasetHandler(svc, processor)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	permRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, permRepo)

	upload := r.App.Group("/api/admin/upload").Use(mdw.AuthMiddleware())
	{
		upload.POST("/:type", mdw.PermissionMiddleware("datasets", "create"), h.Upload)
	}

	ds := r.App.Group("/api/admin/datasets").Use(mdw.AuthMiddleware())
	{
		ds.GET("", mdw.PermissionMiddleware("datasets", "list"), h.List)
		ds.PUT("/:id/status", mdw.PermissionMiddleware("datasets", "update"), h.UpdateStatus)
	}
}

func (r *Routes) PersonRoutes() {
	repo := personRepo.NewPersonRepo(r.DB)
	svc := personSvc.NewPersonService(repo)
	h := personHandler.NewPersonHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	r.App.GET("/api/persons", mdw.AuthMiddleware(), mdw.PermissionMiddleware("persons", "list"), h.GetAll)

	person := r.App.Group("/api/person").Use(mdw.AuthMiddleware())
	{
		person.POST("", mdw.PermissionMiddleware("persons", "create"), h.Create)
		person.GET("/:id", mdw.PermissionMiddleware("persons", "view"), h.GetByID)
		person.PUT("/:id", mdw.PermissionMiddleware("persons", "update"), h.Update)
		person.DELETE("/:id", mdw.PermissionMiddleware("persons", "delete"), h.Delete)
	}
}

func (r *Routes) SessionRoutes() {
	redisClient := database.GetRedisClient()
	if redisClient == nil {
		logger.WriteLog(logger.LogLevelDebug, "Redis not available, session routes will not be registered")
		return
	}

	repo := sessionRepo.NewSessionRepository(redisClient)
	svc := sessionSvc.NewSessionService(repo)
	h := sessionHandler.NewSessionHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// Session management endpoints (authenticated users only)
	sessionGroup := r.App.Group("/api/user").Use(mdw.AuthMiddleware())
	{
		sessionGroup.GET("/sessions", h.GetActiveSessions)
		sessionGroup.DELETE("/session/:session_id", h.RevokeSession)
		sessionGroup.POST("/sessions/revoke-others", h.RevokeAllOtherSessions)
	}

	logger.WriteLog(logger.LogLevelInfo, "Session management routes registered")
}

func (r *Routes) TLRoutes() {
	// Initialize repositories
	mediaRepository := mediaRepo.NewMediaRepo(r.DB)
	activityRepo := tlActivityRepo.NewTLActivityRepo(r.DB)
	attendanceRepo := tlAttendanceRepo.NewTLAttendanceRepo(r.DB)
	sessionRepo := tlSessionRepo.NewTLSessionRepo(r.DB)
	trainingRepo := tlTrainingRepo.NewTLTrainingRepo(r.DB)

	// Initialize services
	mediaService := mediaSvc.NewMediaService(mediaRepository)
	activityService := tlActivitySvc.NewTLActivityService(activityRepo, mediaService)
	attendanceService := tlAttendanceSvc.NewTLAttendanceService(attendanceRepo)
	sessionService := tlSessionSvc.NewTLSessionService(sessionRepo, mediaService)
	trainingService := tlTrainingSvc.NewTLTrainingService(trainingRepo)

	// Initialize handlers
	activityHandler := tlHandler.NewTLActivityHandler(activityService)
	attendanceHandler := tlHandler.NewTLAttendanceHandler(attendanceService)
	sessionHandler := tlHandler.NewTLSessionHandler(sessionService)
	trainingHandler := tlHandler.NewTLTrainingHandler(trainingService)

	// Initialize middleware
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// TL Daily Activity Routes
	activity := r.App.Group("/api/tl/activity").Use(mdw.AuthMiddleware())
	{
		activity.POST("", mdw.PermissionMiddleware("tl_activities", "create"), activityHandler.Create)
		activity.GET("", mdw.PermissionMiddleware("tl_activities", "list"), activityHandler.GetAll)
		activity.GET("/:id", mdw.PermissionMiddleware("tl_activities", "view"), activityHandler.GetByID)
		activity.PUT("/:id", mdw.PermissionMiddleware("tl_activities", "update"), activityHandler.Update)
		activity.DELETE("/:id", mdw.PermissionMiddleware("tl_activities", "delete"), activityHandler.Delete)
	}

	// TL Attendance Routes
	attendance := r.App.Group("/api/tl/attendance").Use(mdw.AuthMiddleware())
	{
		attendance.POST("", mdw.PermissionMiddleware("tl_attendance", "create"), attendanceHandler.Create)
		attendance.GET("", mdw.PermissionMiddleware("tl_attendance", "list"), attendanceHandler.GetAll)
		attendance.GET("/:record_unique_id", mdw.PermissionMiddleware("tl_attendance", "view"), attendanceHandler.GetByRecordUniqueId)
		attendance.PUT("/:record_unique_id", mdw.PermissionMiddleware("tl_attendance", "update"), attendanceHandler.Update)
		attendance.DELETE("/:record_unique_id", mdw.PermissionMiddleware("tl_attendance", "delete"), attendanceHandler.Delete)
	}

	// TL Session Routes (Merged Coaching & Briefing)
	session := r.App.Group("/api/tl/session").Use(mdw.AuthMiddleware())
	{
		session.POST("", mdw.PermissionMiddleware("tl_sessions", "create"), sessionHandler.Create)
		session.GET("", mdw.PermissionMiddleware("tl_sessions", "list"), sessionHandler.GetAll)
		session.GET("/:id", mdw.PermissionMiddleware("tl_sessions", "view"), sessionHandler.GetByID)
		session.PUT("/:id", mdw.PermissionMiddleware("tl_sessions", "update"), sessionHandler.Update)
		session.DELETE("/:id", mdw.PermissionMiddleware("tl_sessions", "delete"), sessionHandler.Delete)
	}

	// TL Training Participation Routes
	training := r.App.Group("/api/tl/training").Use(mdw.AuthMiddleware())
	{
		training.POST("", mdw.PermissionMiddleware("tl_training", "create"), trainingHandler.Create)
		training.GET("", mdw.PermissionMiddleware("tl_training", "list"), trainingHandler.GetAll)
		training.GET("/:training_batch", mdw.PermissionMiddleware("tl_training", "view"), trainingHandler.GetByTrainingBatch)
	}

	logger.WriteLog(logger.LogLevelInfo, "Team Leader input routes registered")
}
