package handler

import (
	"backend/internal/api/dto"
	"backend/internal/api/vo"
	"backend/internal/model"
	promptService "backend/internal/service/prompt"
	versionService "backend/internal/service/version"
	"backend/pkg/errors"
	"backend/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type PromptHandler struct {
	service        *promptService.Service
	versionService *versionService.Service
}

func CreatePromptHandler(service *promptService.Service, versionService *versionService.Service) *PromptHandler {
	return &PromptHandler{
		service:        service,
		versionService: versionService,
	}
}

func (s *PromptHandler) Create(c *gin.Context) {
	var req dto.CreatePromptDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}
	if strings.TrimSpace(req.Path) == "" {
		req.Path = fmt.Sprintf("/%s", req.Name)
	}
	p, err := s.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: "internal server error",
		})
		return
	}

	response.Success(c, vo.FromPrompt(p))
}

func (s *PromptHandler) GetPromptByID(c *gin.Context) {
	promptId := c.Param("id")
	p, err := s.service.GetByID(c.Request.Context(), promptId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.FromPrompt(p))
}

func (s *PromptHandler) GetPromptByPath(c *gin.Context) {
	path := c.Param("path")
	if strings.TrimSpace(path) == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}
	p, err := s.service.GetByPath(c.Request.Context(), path)
	if p == nil && err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	// 如果未发布，返回无发布版本
	if p != nil && !p.IsPublish {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "no published version",
		})
		return
	}

	// 根据prompt表的version字段(存储的是版本ID)查询版本详情
	version, err := s.versionService.GetByID(c.Request.Context(), p.LatestVersion)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	response.Success(c, gin.H{
		"prompt":  vo.FromPrompt(p),
		"version": vo.FromPromptVersion(version),
	})
}

func (s *PromptHandler) Update(c *gin.Context) {
	var req dto.UpdatePromptDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid request body",
		})
		return
	}

	p := &model.Prompt{
		ID:        req.ID,
		Name:      req.Name,
		IsPublish: req.IsPublish,
		Category:  req.Category,
	}

	if err := s.service.Update(c.Request.Context(), p); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (s *PromptHandler) Delete(c *gin.Context) {
	promptId := c.Param("id")
	if promptId == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "invalid prompt id",
		})
		return
	}

	if err := s.service.DeleteByID(c.Request.Context(), promptId); err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, nil)
}

func (s *PromptHandler) List(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    errors.DefaultError,
			Data:    nil,
			Message: "userId is required",
		})
		return
	}

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

	prompts, total, err := s.service.List(c.Request.Context(), username, offset, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.Response{
			Code:    errors.ServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	response.Success(c, vo.NewPageData(vo.FromPrompts(prompts), total, offset, limit))
}

func (s *PromptHandler) ReverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 覆盖 Director，固定转发路径
	proxy.Director = func(req *http.Request) {
		// 保留原始 method / query / body
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.URL.Path = targetURL.Path
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
