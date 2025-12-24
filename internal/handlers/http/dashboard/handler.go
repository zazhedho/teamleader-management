package handlerdashboard

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"teamleader-management/internal/dto"
	servicedashboard "teamleader-management/internal/services/dashboard"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	TLDashboardService    *servicedashboard.TLDashboardService
	AdminAnalyticsService *servicedashboard.AdminAnalyticsService
}

func NewDashboardHandler(
	tlDashboardService *servicedashboard.TLDashboardService,
	adminAnalyticsService *servicedashboard.AdminAnalyticsService,
) *DashboardHandler {
	return &DashboardHandler{
		TLDashboardService:    tlDashboardService,
		AdminAnalyticsService: adminAnalyticsService,
	}
}

// GetMyDashboard retrieves dashboard for the currently authenticated TL
// GET /api/dashboard/me
func (h *DashboardHandler) GetMyDashboard(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DashboardHandler][GetMyDashboard]", logId)

	// Get authenticated user
	authUser, exists := ctx.Get("user")
	if !exists {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; User not authenticated", logPrefix))
		res := response.Response(http.StatusUnauthorized, "Unauthorized", logId, nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	user := authUser.(map[string]interface{})
	personId := user["person_id"].(string)

	// Get period from query params or use current month/year
	periodMonth := ctx.DefaultQuery("period_month", strconv.Itoa(int(time.Now().Month())))
	periodYear := ctx.DefaultQuery("period_year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(periodMonth)
	if err != nil || month < 1 || month > 12 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_month: %s", logPrefix, periodMonth))
		res := response.Response(http.StatusBadRequest, "Invalid period_month", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	year, err := strconv.Atoi(periodYear)
	if err != nil || year < 2020 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_year: %s", logPrefix, periodYear))
		res := response.Response(http.StatusBadRequest, "Invalid period_year", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get dashboard
	dashboard, err := h.TLDashboardService.GetDashboard(personId, month, year)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetDashboard; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, "Failed to get dashboard: "+err.Error(), logId, nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Dashboard retrieved successfully", logId, dashboard)
	ctx.JSON(http.StatusOK, res)
}

// GetDashboardByPersonId retrieves dashboard for a specific TL (admin only)
// GET /api/dashboard/:person_id
func (h *DashboardHandler) GetDashboardByPersonId(ctx *gin.Context) {
	personId := ctx.Param("person_id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DashboardHandler][GetDashboardByPersonId]", logId)

	if personId == "" {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; person_id is required", logPrefix))
		res := response.Response(http.StatusBadRequest, "person_id is required", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get period from query params or use current month/year
	periodMonth := ctx.DefaultQuery("period_month", strconv.Itoa(int(time.Now().Month())))
	periodYear := ctx.DefaultQuery("period_year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(periodMonth)
	if err != nil || month < 1 || month > 12 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_month: %s", logPrefix, periodMonth))
		res := response.Response(http.StatusBadRequest, "Invalid period_month", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	year, err := strconv.Atoi(periodYear)
	if err != nil || year < 2020 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_year: %s", logPrefix, periodYear))
		res := response.Response(http.StatusBadRequest, "Invalid period_year", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get dashboard
	dashboard, err := h.TLDashboardService.GetDashboard(personId, month, year)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetDashboard; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, "Failed to get dashboard: "+err.Error(), logId, nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Dashboard retrieved successfully", logId, dashboard)
	ctx.JSON(http.StatusOK, res)
}

// GetAnalytics retrieves comprehensive analytics for admin
// GET /api/analytics
func (h *DashboardHandler) GetAnalytics(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DashboardHandler][GetAnalytics]", logId)

	// Get period from query params or use current month/year
	periodMonth := ctx.DefaultQuery("period_month", strconv.Itoa(int(time.Now().Month())))
	periodYear := ctx.DefaultQuery("period_year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(periodMonth)
	if err != nil || month < 1 || month > 12 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_month: %s", logPrefix, periodMonth))
		res := response.Response(http.StatusBadRequest, "Invalid period_month", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	year, err := strconv.Atoi(periodYear)
	if err != nil || year < 2020 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_year: %s", logPrefix, periodYear))
		res := response.Response(http.StatusBadRequest, "Invalid period_year", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get optional parameters
	topN := 5 // default top/bottom performers count
	if topNParam := ctx.Query("top_n"); topNParam != "" {
		if n, err := strconv.Atoi(topNParam); err == nil && n > 0 && n <= 20 {
			topN = n
		}
	}

	trendMonths := 6 // default trend comparison months
	if trendParam := ctx.Query("trend_months"); trendParam != "" {
		if n, err := strconv.Atoi(trendParam); err == nil && n > 0 && n <= 12 {
			trendMonths = n
		}
	}

	// Get analytics
	analytics, err := h.AdminAnalyticsService.GetAnalytics(month, year, topN, trendMonths)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetAnalytics; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, "Failed to get analytics: "+err.Error(), logId, nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Analytics retrieved successfully", logId, analytics)
	ctx.JSON(http.StatusOK, res)
}

// CompareTeam compares multiple TLs side by side
// POST /api/analytics/compare
func (h *DashboardHandler) CompareTeam(ctx *gin.Context) {
	var req dto.TeamComparisonRequest
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DashboardHandler][CompareTeam]", logId)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate person_ids count
	if len(req.PersonIds) < 2 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; At least 2 person IDs are required", logPrefix))
		res := response.Response(http.StatusBadRequest, "At least 2 person IDs are required", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if len(req.PersonIds) > 10 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Maximum 10 person IDs allowed", logPrefix))
		res := response.Response(http.StatusBadRequest, "Maximum 10 person IDs allowed", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate period
	if req.PeriodMonth < 1 || req.PeriodMonth > 12 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_month: %d", logPrefix, req.PeriodMonth))
		res := response.Response(http.StatusBadRequest, "Invalid period_month", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if req.PeriodYear < 2020 {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Invalid period_year: %d", logPrefix, req.PeriodYear))
		res := response.Response(http.StatusBadRequest, "Invalid period_year", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Compare team
	comparison, err := h.AdminAnalyticsService.CompareTeam(req.PersonIds, req.PeriodMonth, req.PeriodYear)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.CompareTeam; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, "Failed to compare team: "+err.Error(), logId, nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Team comparison retrieved successfully", logId, comparison)
	ctx.JSON(http.StatusOK, res)
}
