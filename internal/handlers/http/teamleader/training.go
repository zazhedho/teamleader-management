package handlerteamleader

import (
	"fmt"
	"net/http"
	"reflect"
	"teamleader-management/internal/dto"
	interfacetltraining "teamleader-management/internal/interfaces/tltraining"
	"teamleader-management/pkg/filter"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type TLTrainingHandler struct {
	Service interfacetltraining.ServiceTLTrainingInterface
}

func NewTLTrainingHandler(s interfacetltraining.ServiceTLTrainingInterface) *TLTrainingHandler {
	return &TLTrainingHandler{Service: s}
}

func (h *TLTrainingHandler) Create(ctx *gin.Context) {
	var req dto.TLTrainingCreate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLTrainingHandler][Create]", logId)

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

	res := response.Response(http.StatusCreated, "Training participation created successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *TLTrainingHandler) GetByTrainingBatch(ctx *gin.Context) {
	trainingBatch := ctx.Param("training_batch")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLTrainingHandler][GetByTrainingBatch]", logId)

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

	data, err := h.Service.GetByTrainingBatch(trainingBatch, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetByTrainingBatch; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Training record not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get training record successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLTrainingHandler) GetAll(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLTrainingHandler][GetAll]", logId)

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
