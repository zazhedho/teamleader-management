package handlerteamleader

import (
	"fmt"
	"net/http"
	"reflect"
	"teamleader-management/internal/dto"
	interfacemedia "teamleader-management/internal/interfaces/media"
	interfacetlsession "teamleader-management/internal/interfaces/tlsession"
	"teamleader-management/pkg/filter"
	"teamleader-management/pkg/logger"
	"teamleader-management/pkg/messages"
	"teamleader-management/pkg/response"
	"teamleader-management/utils"

	"github.com/gin-gonic/gin"
)

type TLSessionHandler struct {
	Service      interfacetlsession.ServiceTLSessionInterface
	MediaService interfacemedia.ServiceMediaInterface
}

func NewTLSessionHandler(s interfacetlsession.ServiceTLSessionInterface, mediaService interfacemedia.ServiceMediaInterface) *TLSessionHandler {
	return &TLSessionHandler{
		Service:      s,
		MediaService: mediaService,
	}
}

func (h *TLSessionHandler) Create(ctx *gin.Context) {
	var req dto.TLSessionCreate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][Create]", logId)

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

	sessionType := req.SessionType
	message := fmt.Sprintf("%s session created successfully", sessionType)
	res := response.Response(http.StatusCreated, message, logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *TLSessionHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][GetByID]", logId)

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
		res := response.Response(http.StatusNotFound, "Session not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get session successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLSessionHandler) GetAll(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][GetAll]", logId)

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

func (h *TLSessionHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.TLSessionUpdate
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][Update]", logId)

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

	res := response.Response(http.StatusOK, "Session updated successfully", logId, data)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLSessionHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][Delete]", logId)

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

	res := response.Response(http.StatusOK, "Session deleted successfully", logId, nil)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("%s; Success", logPrefix))
	ctx.JSON(http.StatusOK, res)
}

func (h *TLSessionHandler) UploadFile(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][UploadFile]", logId)

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

	// Verify session ownership
	_, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Session not found or unauthorized: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Session not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Check existing media count (max 2)
	existingMedia, _ := h.MediaService.GetMediaByEntity(utils.EntityTLSession, id)
	if len(existingMedia) >= 2 {
		res := response.Response(http.StatusBadRequest, "Maximum 2 photos allowed per session", logId, nil)
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

	data, err := h.MediaService.UploadAndAttach(ctx.Request.Context(), utils.EntityTLSession, id, file, actor)
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

func (h *TLSessionHandler) GetFiles(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][GetFiles]", logId)

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

	// Verify session ownership
	_, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Session not found or unauthorized: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Session not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	data, err := h.MediaService.GetMediaByEntity(utils.EntityTLSession, id)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; MediaService.GetMediaByEntity; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Get session files successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *TLSessionHandler) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")
	fileId := ctx.Param("fileId")
	logId := utils.GenerateLogId(ctx)
	logPrefix := fmt.Sprintf("[%s][TLSessionHandler][DeleteFile]", logId)

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

	// Verify session ownership
	_, err := h.Service.GetByID(id, personId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; Session not found or unauthorized: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Session not found", logId, nil)
		res.Error = err.Error()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Verify file belongs to this session
	media, err := h.MediaService.GetMediaByID(fileId)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("%s; File not found: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "File not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if media.EntityType != utils.EntityTLSession || media.EntityId != id {
		res := response.Response(http.StatusForbidden, "File does not belong to this session", logId, nil)
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
