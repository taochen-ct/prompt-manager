package handler

import (
	"backend/internal/api/dto"
	"backend/internal/api/middleware"
	"backend/internal/api/vo"
	favoritesService "backend/internal/service/favorites"
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteHandler struct {
	service *favoritesService.Service
}

func CreateFavoriteHandler(service *favoritesService.Service) *FavoriteHandler {
	return &FavoriteHandler{
		service: service,
	}
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	var req dto.AddFavoriteDTO
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

	f, err := h.service.AddFavorite(c.Request.Context(), userID, req)
	if err != nil {
		if err == favoritesService.ErrFavoriteExists {
			response.Error(c, http.StatusConflict, response.Response{
				Code:    errors.DefaultError,
				Data:    nil,
				Message: "already in favorites",
			})
			return
		}
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	response.Success(c, vo.FromFavorite(f))
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	var req dto.RemoveFavoriteDTO
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

	if err := h.service.RemoveFavorite(c.Request.Context(), userID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (h *FavoriteHandler) Check(c *gin.Context) {
	var req dto.CheckFavoriteDTO
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

	isFavorite, err := h.service.IsFavorite(c.Request.Context(), userID, req.PromptID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	response.Success(c, gin.H{"isFavorite": isFavorite})
}

func (h *FavoriteHandler) List(c *gin.Context) {
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

	list, total, err := h.service.ListFavorites(c.Request.Context(), userID, offset, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.NewPageData(vo.FromFavoritesWithPrompt(list), total, offset, limit))
}
