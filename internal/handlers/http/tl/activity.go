package handlertl

import (
	"fmt"
	"net/http"
	"reflect"
	"teamleader-management/internal/dto"
	interfacemedia "teamleader-management/internal/interfaces/media"
	interfacetlactivity "teamleader-management/internal/interfaces/tlactivity"
	"teamleader-management/pkg/filter"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type TLActivityHandler struct {
	Service      interfacetlactivity.ServiceTLActivityInterface
	MediaService interfacemedia.ServiceMediaInterface
}

func NewTLActivityHandler(s interfacetlactivity.ServiceTLActivityInterface, mediaService interfacemedia.ServiceMediaInterface) *TLActivityHandler {
	return &TLActivityHandler{
		Service:      s,
		MediaService: mediaService,
	}
}

func (h *TLActivityHandler) Create(ctx *gin.Context) {
	var req dto.TLActivityCreate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][Create]", logId)

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

	res := response.Response(http.StatusCreated, "Daily activity created successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *TLActivityHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][GetByID]", logId)

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

	data, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.GetByID; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Daily activity not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get daily activity successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLActivityHandler) GetAll(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][GetAll]", logId)

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

func (h *TLActivityHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.TLActivityUpdate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][Update]", logId)

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

	data, err := h.Service.Update(id, personId, req, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Update; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusOK, "Daily activity updated successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLActivityHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][Delete]", logId)

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

	if err := h.Service.Delete(ctx.Request.Context(), id, personId); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Service.Delete; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, err.Error(), logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Daily activity deleted successfully", logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success", logPrefix))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLActivityHandler) UploadFile(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][UploadFile]", logId)

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

	// Verify activity ownership
	_, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Activity not found or unauthorized: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Activity not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Check existing media count (max 2)
	existingMedia, _ := h.MediaService.GetMediaByEntity("tl_activity", id)
	if len(existingMedia) >= 2 {
		res := response.Response(http.StatusBadRequest, "Maximum 2 photos allowed per activity", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; FormFile ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = "file is required"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.MediaService.UploadAndAttach(ctx.Request.Context(), "tl_activity", id, file, actor)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; MediaService.UploadAndAttach; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusCreated, "File uploaded successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *TLActivityHandler) GetFiles(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][GetFiles]", logId)

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

	// Verify activity ownership
	_, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Activity not found or unauthorized: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Activity not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	data, err := h.MediaService.GetMediaByEntity("tl_activity", id)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; MediaService.GetMediaByEntity; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Get activity files successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *TLActivityHandler) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")
	fileId := ctx.Param("fileId")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLActivityHandler][DeleteFile]", logId)

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

	// Verify activity ownership
	_, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Activity not found or unauthorized: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Activity not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Verify file belongs to this activity
	media, err := h.MediaService.GetMediaByID(fileId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; File not found: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "File not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if media.EntityType != "tl_activity" || media.EntityId != id {
		res := response.Response(http.StatusForbidden, "File does not belong to this activity", logId, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	if err := h.MediaService.DeleteMediaByID(ctx.Request.Context(), fileId); err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; MediaService.DeleteMediaByID; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "File deleted successfully", logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success", logPrefix))
	ctx.JSON(http.StatusOK, res)
}
