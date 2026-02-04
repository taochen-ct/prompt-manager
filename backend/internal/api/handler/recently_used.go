package handler

import (
	"backend/internal/api/dto"
	"backend/internal/api/middleware"
	"backend/internal/api/vo"
	recentlyUsedService "backend/internal/service/recently_used"
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RecentlyUsedHandler struct {
	service *recentlyUsedService.Service
}

func CreateRecentlyUsedHandler(service *recentlyUsedService.Service) *RecentlyUsedHandler {
	return &RecentlyUsedHandler{
		service: service,
	}
}

func (h *RecentlyUsedHandler) Record(c *gin.Context) {
	var req dto.RecordUsageDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	userID, _, ok := middleware.GetUserFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "unauthorized",
		})
		return
	}

	rec, err := h.service.RecordUsage(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	response.Success(c, vo.FromRecentlyUsed(rec))
}

func (h *RecentlyUsedHandler) Remove(c *gin.Context) {
	var req dto.RemoveRecentlyUsedDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	userID, _, ok := middleware.GetUserFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "unauthorized",
		})
		return
	}

	if err := h.service.RemoveRecord(c.Request.Context(), userID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (h *RecentlyUsedHandler) List(c *gin.Context) {
	userID, _, ok := middleware.GetUserFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "unauthorized",
		})
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	list, total, err := h.service.ListRecentlyUsed(c.Request.Context(), userID, offset, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.NewPageData(vo.FromRecentlyUsedWithPromptList(list), total, offset, limit))
}

func (h *RecentlyUsedHandler) Clean(c *gin.Context) {
	var req dto.CleanRecentlyUsedDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	userID, _, ok := middleware.GetUserFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "unauthorized",
		})
		return
	}

	if err := h.service.CleanOldRecords(c.Request.Context(), userID, req.KeepCount); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}
