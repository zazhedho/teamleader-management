package handlerdataset

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"teamleader-management/internal/dto"
	interfacedataset "teamleader-management/internal/interfaces/dataset"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type DatasetHandler struct {
	Service   interfacedataset.ServiceDatasetInterface
	Processor interfacedataset.DatasetProcessorInterface
}

func NewDatasetHandler(s interfacedataset.ServiceDatasetInterface, p interfacedataset.DatasetProcessorInterface) *DatasetHandler {
	return &DatasetHandler{Service: s, Processor: p}
}

func (h *DatasetHandler) Upload(ctx *gin.Context) {
	datasetType := ctx.Param("type")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DatasetHandler][Upload]", logId)

	periodDate := ctx.PostForm("period_date")
	periodMonthStr := ctx.PostForm("period_month")
	periodYearStr := ctx.PostForm("period_year")

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; FormFile ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = "file is required"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	defer file.Close()

	authData := utils.GetAuthData(ctx)
	actor := ""
	if authData != nil {
		actor = utils.InterfaceString(authData["user_id"])
	}

	req := dto.DatasetUploadRequest{
		PeriodDate: periodDate,
		Type:       datasetType,
	}
	if periodMonthStr != "" {
		if num, err := strconv.Atoi(periodMonthStr); err == nil {
			req.PeriodMonth = num
		}
	}
	if periodYearStr != "" {
		if num, err := strconv.Atoi(periodYearStr); err == nil {
			req.PeriodYear = num
		}
	}
	if v := ctx.PostForm("period_frequency"); v != "" {
		req.PeriodFrequency = v
	}

	ds, data, err := h.Service.Create(datasetType, req, file, fileHeader, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Upload; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	processed, err := h.Processor.ProcessStream(&ds, data, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Processor.Process; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusCreated, "Dataset uploaded and processed successfully", logId, processed)
	ctx.JSON(http.StatusCreated, res)
}

func (h *DatasetHandler) List(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DatasetHandler][List]", logId)

	filters := map[string]interface{}{}
	if v := ctx.Query("type"); v != "" {
		filters["type"] = strings.ToUpper(v)
	}
	if v := ctx.Query("status"); v != "" {
		filters["status"] = strings.ToUpper(v)
	}
	if v := ctx.Query("period_year"); v != "" {
		if num, err := strconv.Atoi(v); err == nil {
			filters["period_year"] = num
		}
	}
	if v := ctx.Query("period_month"); v != "" {
		if num, err := strconv.Atoi(v); err == nil {
			filters["period_month"] = num
		}
	}

	data, total, err := h.Service.List(filters)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.List; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.PaginationResponse(http.StatusOK, int(total), 1, len(data), logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *DatasetHandler) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][DatasetHandler][UpdateStatus]", logId)

	var req dto.DatasetStatusUpdate
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	authData := utils.GetAuthData(ctx)
	actor := ""
	if authData != nil {
		actor = utils.InterfaceString(authData["user_id"])
	}

	data, err := h.Service.UpdateStatus(id, req.Status, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.UpdateStatus; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusOK, "Dataset status updated", logId, data)
	ctx.JSON(http.StatusOK, res)
}
