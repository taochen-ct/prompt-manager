package handler

import (
	"backend/internal/api/dto"
	"backend/internal/api/vo"
	"backend/internal/model"
	versionService "backend/internal/service/version"
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PromptVersionHandler struct {
	service *versionService.Service
}

func CreatePromptVersionHandler(service *versionService.Service) *PromptVersionHandler {
	return &PromptVersionHandler{
		service: service,
	}
}

func (h *PromptVersionHandler) Create(c *gin.Context) {
	var req dto.CreatePromptVersionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	v, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromPromptVersion(v))
}

func (h *PromptVersionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid version id",
		})
		return
	}

	v, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromPromptVersion(v))
}

func (h *PromptVersionHandler) GetByPromptID(c *gin.Context) {
	promptID := c.Param("promptId")
	if promptID == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid prompt id",
		})
		return
	}

	versions, err := h.service.GetByPromptID(c.Request.Context(), promptID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromPromptVersions(versions))
}

func (h *PromptVersionHandler) GetLatestByPromptID(c *gin.Context) {
	promptID := c.Param("promptId")
	if promptID == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid prompt id",
		})
		return
	}

	v, err := h.service.GetLatestByPromptID(c.Request.Context(), promptID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromPromptVersion(v))
}

func (h *PromptVersionHandler) Update(c *gin.Context) {
	var req dto.UpdatePromptVersionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	v := &model.PromptVersion{
		ID:        req.ID,
		Version:   req.Version,
		Content:   req.Content,
		Variables: req.Variables,
		ChangeLog: req.ChangeLog,
		IsPublish: req.IsPublish,
	}

	if err := h.service.Update(c.Request.Context(), v); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (h *PromptVersionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid version id",
		})
		return
	}

	if err := h.service.DeleteByID(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (h *PromptVersionHandler) List(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	versions, total, err := h.service.List(c.Request.Context(), offset, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.NewPageData(vo.FromPromptVersions(versions), total, offset, limit))
}
