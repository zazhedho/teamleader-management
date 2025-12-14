package handlerteamleader

import (
	"fmt"
	"net/http"
	"reflect"
	"teamleader-management/internal/dto"
	interfacetlattendance "teamleader-management/internal/interfaces/tlattendance"
	"teamleader-management/pkg/filter"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type TLAttendanceHandler struct {
	Service interfacetlattendance.ServiceTLAttendanceInterface
}

func NewTLAttendanceHandler(s interfacetlattendance.ServiceTLAttendanceInterface) *TLAttendanceHandler {
	return &TLAttendanceHandler{Service: s}
}

func (h *TLAttendanceHandler) Create(ctx *gin.Context) {
	var req dto.TLAttendanceCreate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLAttendanceHandler][Create]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	authData := utils.GetAuthData(ctx)
	if authData == nil {
		res := response.Response(http.StatusUnauthorized, "Unauthorized", logId, nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	actor := utils.InterfaceString(authData["user_id"])
	personId := utils.InterfaceString(authData["person_id"])

	if personId == "" {
		res := response.Response(http.StatusForbidden, "User is not linked to a person", logId, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	data, err := h.Service.Create(personId, req, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Create; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusCreated, "Attendance created successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *TLAttendanceHandler) GetByRecordUniqueId(ctx *gin.Context) {
	recordUniqueId := ctx.Param("record_unique_id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLAttendanceHandler][GetByRecordUniqueId]", logId)

	authData := utils.GetAuthData(ctx)
	if authData == nil {
		res := response.Response(http.StatusUnauthorized, "Unauthorized", logId, nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	personId := utils.InterfaceString(authData["person_id"])
	if personId == "" {
		res := response.Response(http.StatusForbidden, "User is not linked to a person", logId, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	data, err := h.Service.GetByRecordUniqueId(recordUniqueId, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetByRecordUniqueId; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Attendance record not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get attendance record successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLAttendanceHandler) GetAll(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLAttendanceHandler][GetAll]", logId)

	authData := utils.GetAuthData(ctx)
	if authData == nil {
		res := response.Response(http.StatusUnauthorized, "Unauthorized", logId, nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	personId := utils.InterfaceString(authData["person_id"])
	if personId == "" {
		res := response.Response(http.StatusForbidden, "User is not linked to a person", logId, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	params, err := filter.GetBaseParams(ctx, "date", "desc", 10)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; GetBaseParams; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, total, err := h.Service.GetAll(personId, params)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetAll; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.PaginationResponse(http.StatusOK, int(total), params.Page, params.Limit, logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLAttendanceHandler) Update(ctx *gin.Context) {
	recordUniqueId := ctx.Param("record_unique_id")
	var req dto.TLAttendanceUpdate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLAttendanceHandler][Update]", logId)

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	authData := utils.GetAuthData(ctx)
	if authData == nil {
		res := response.Response(http.StatusUnauthorized, "Unauthorized", logId, nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	actor := utils.InterfaceString(authData["user_id"])
	personId := utils.InterfaceString(authData["person_id"])

	if personId == "" {
		res := response.Response(http.StatusForbidden, "User is not linked to a person", logId, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	data, err := h.Service.Update(recordUniqueId, personId, req, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Update; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusOK, "Attendance updated successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLAttendanceHandler) Delete(ctx *gin.Context) {
	recordUniqueId := ctx.Param("record_unique_id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLAttendanceHandler][Delete]", logId)

	authData := utils.GetAuthData(ctx)
	if authData == nil {
		res := response.Response(http.StatusUnauthorized, "Unauthorized", logId, nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	personId := utils.InterfaceString(authData["person_id"])
	if personId == "" {
		res := response.Response(http.StatusForbidden, "User is not linked to a person", logId, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	if err := h.Service.Delete(recordUniqueId, personId); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Delete; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, err.Error(), logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Attendance deleted successfully", logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success", logPrefix))
	ctx.JSON(http.StatusOK, res)
}
