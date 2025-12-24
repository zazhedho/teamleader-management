package handlerevaluation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"teamleader-management/internal/dto"
	interfaceevaluation "teamleader-management/internal/interfaces/evaluation"
	"teamleader-management/pkg/filter"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type EvaluationHandler struct {
	Service interfaceevaluation.ServiceEvaluationInterface
}

func NewEvaluationHandler(s interfaceevaluation.ServiceEvaluationInterface) *EvaluationHandler {
	return &EvaluationHandler{Service: s}
}

// CalculateEvaluation calculates evaluation for a period
// POST /api/evaluation/calculate
func (h *EvaluationHandler) CalculateEvaluation(ctx *gin.Context) {
	var req dto.EvaluationCalculateRequest
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][EvaluationHandler][CalculateEvaluation]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	// Calculate evaluation
	results, err := h.Service.CalculateEvaluation(req.PeriodMonth, req.PeriodYear, req.PersonId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.CalculateEvaluation; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	message := fmt.Sprintf("Calculated evaluation for %d TL(s) in period %d-%02d", len(results), req.PeriodYear, req.PeriodMonth)
	res := response.Response(http.StatusOK, message, logId, results)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success; Count: %d", logPrefix, len(results)))
	ctx.JSON(http.StatusOK, res)
}

// RecalculateEvaluation recalculates evaluation (deletes old and creates new)
// POST /api/evaluation/recalculate
func (h *EvaluationHandler) RecalculateEvaluation(ctx *gin.Context) {
	var req dto.EvaluationCalculateRequest
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][EvaluationHandler][RecalculateEvaluation]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	results, err := h.Service.RecalculateEvaluation(req.PeriodMonth, req.PeriodYear, req.PersonId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.RecalculateEvaluation; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	message := fmt.Sprintf("Recalculated evaluation for %d TL(s)", len(results))
	res := response.Response(http.StatusOK, message, logId, results)
	ctx.JSON(http.StatusOK, res)
}

// GetByID retrieves evaluation by ID with full breakdown
// GET /api/evaluation/:id
func (h *EvaluationHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][EvaluationHandler][GetByID]", logId)

	data, err := h.Service.GetByID(id)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetByID; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Evaluation not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get evaluation successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success", logPrefix))
	ctx.JSON(http.StatusOK, res)
}

// GetByPersonAndPeriod retrieves evaluation for a person in a specific period
// GET /api/evaluation/person/:person_id
// Query params: period_month, period_year
func (h *EvaluationHandler) GetByPersonAndPeriod(ctx *gin.Context) {
	personId := ctx.Param("person_id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][EvaluationHandler][GetByPersonAndPeriod]", logId)

	periodMonthStr := ctx.Query("period_month")
	periodYearStr := ctx.Query("period_year")

	if periodMonthStr == "" || periodYearStr == "" {
		res := response.Response(http.StatusBadRequest, "period_month and period_year are required", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	periodMonth, err := strconv.Atoi(periodMonthStr)
	if err != nil || periodMonth < 1 || periodMonth > 12 {
		res := response.Response(http.StatusBadRequest, "invalid period_month", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	periodYear, err := strconv.Atoi(periodYearStr)
	if err != nil || periodYear < 2020 {
		res := response.Response(http.StatusBadRequest, "invalid period_year", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.GetByPersonAndPeriod(personId, periodMonth, periodYear)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetByPersonAndPeriod; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Evaluation not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get evaluation successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

// GetAll lists evaluations with pagination and filters
// GET /api/evaluation
func (h *EvaluationHandler) GetAll(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][EvaluationHandler][GetAll]", logId)

	params, err := filter.GetBaseParams(ctx, "created_at", "desc", 10)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; GetBaseParams; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, total, err := h.Service.GetAll(params)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetAll; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.PaginationResponse(http.StatusOK, int(total), params.Page, params.Limit, logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success; Total: %d", logPrefix, total))
	ctx.JSON(http.StatusOK, res)
}

// GetLeaderboard retrieves top TLs for a period
// GET /api/evaluation/leaderboard
// Query params: period_month, period_year, limit (optional)
func (h *EvaluationHandler) GetLeaderboard(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][EvaluationHandler][GetLeaderboard]", logId)

	periodMonthStr := ctx.Query("period_month")
	periodYearStr := ctx.Query("period_year")
	limitStr := ctx.DefaultQuery("limit", "10")

	if periodMonthStr == "" || periodYearStr == "" {
		res := response.Response(http.StatusBadRequest, "period_month and period_year are required", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	periodMonth, err := strconv.Atoi(periodMonthStr)
	if err != nil || periodMonth < 1 || periodMonth > 12 {
		res := response.Response(http.StatusBadRequest, "invalid period_month", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	periodYear, err := strconv.Atoi(periodYearStr)
	if err != nil || periodYear < 2020 {
		res := response.Response(http.StatusBadRequest, "invalid period_year", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	data, err := h.Service.GetLeaderboard(periodMonth, periodYear, limit)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetLeaderboard; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Get leaderboard successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success; Period: %s", logPrefix, data.Period))
	ctx.JSON(http.StatusOK, res)
}
