package handler

import (
	"backend/internal/api/dto"
	"backend/internal/api/vo"
	"backend/internal/service/user"
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *user.Service
}

func CreateUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	userData, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	response.Success(c, vo.FromUser(userData))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}
	userData, err := h.service.Get(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromUser(userData))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var req dto.DeleteUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}
	userId, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}
	if err = h.service.Delete(c, userId); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
	}
	response.Success(c, nil)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid user id",
		})
		return
	}

	var req dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), userId, req); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}
