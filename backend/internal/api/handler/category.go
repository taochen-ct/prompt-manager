package handler

import (
	"backend/internal/api/dto"
	"backend/internal/api/vo"
	"backend/internal/model"
	categoryService "backend/internal/service/category"
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryHandler struct {
	service *categoryService.Service
}

func CreateCategoryHandler(service *categoryService.Service) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	category, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	response.Success(c, vo.FromCategory(category))
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid category id",
		})
		return
	}

	category, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	if category == nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "category not found",
		})
		return
	}

	response.Success(c, vo.FromCategory(category))
}

func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromCategories(list))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	var req dto.UpdateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	category := &model.Category{
		ID:    req.ID,
		Title: req.Title,
		Icon:  req.Icon,
		Count: req.Count,
		URL:   req.URL,
	}

	if err := h.service.Update(c.Request.Context(), category); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid category id",
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
