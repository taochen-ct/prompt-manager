package router

import (
	"backend/internal/api/handler"
	"backend/internal/api/middleware"
	"backend/pkg/common"
	"backend/pkg/config"
	"backend/pkg/response"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	loggerMiddleware *middleware.Logger,
	recoveryMiddleware *middleware.Recovery,
	corsMiddleware *middleware.Cors,
	jwtMiddleware *middleware.JWTMiddleware,
	userHandler *handler.UserHandler,
	promptHandler *handler.PromptHandler,
	versionHandler *handler.PromptVersionHandler,
	categoryHandler *handler.CategoryHandler,
	favoriteHandler *handler.FavoriteHandler,
	recentlyUsedHandler *handler.RecentlyUsedHandler,
	remoteLogHandler *handler.RemoteLogHandler,
) *gin.Engine {
	if cfg.Server.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(
		gin.Logger(),
		loggerMiddleware.Handler(),
		corsMiddleware.Handler(),
		recoveryMiddleware.Handler(),
	)
	// website server
	if common.IsExist(cfg.Web.StaticDir) {
		r.Use(static.Serve("/", static.LocalFile(cfg.Web.StaticDir, false)))
		r.NoRoute(func(c *gin.Context) {
			c.File(fmt.Sprintf("%s/%s", cfg.Web.StaticDir, cfg.Web.DefaultHtml))
		})
	}

	// base api
	api := r.Group(cfg.Server.ApiPrefix)
	{
		api.GET("/ping", func(c *gin.Context) {
			response.Success(c, "pong")
		})
		api.GET("/prompt/content/*path", promptHandler.GetPromptByPath)
		api.POST("/remote/log/push", remoteLogHandler.Handler)
	}

	// user api collections (公开接口不需要 JWT)
	userAPI := api.Group("/user")
	{
		userAPI.POST("/login", userHandler.Login)
		userAPI.POST("/logout", userHandler.Logout)
		userAPI.POST("/create", userHandler.CreateUser)
		userAPI.POST("/delete", userHandler.DeleteUser)
		userAPI.GET("/info/:id", userHandler.GetUser)
		userAPI.POST("/update/:id", userHandler.UpdateUser)
	}

	// 需要 JWT 认证的路由
	authAPI := api.Group("")
	authAPI.Use(jwtMiddleware.Handler())
	{
		// prompt api
		promptAPI := authAPI.Group("/prompt")
		{
			promptAPI.POST("/create", promptHandler.Create)
			promptAPI.GET("/info/:id", promptHandler.GetPromptByID)
			promptAPI.POST("/update", promptHandler.Update)
			promptAPI.POST("/delete/:id", promptHandler.Delete)
			promptAPI.GET("/list", promptHandler.List)
			promptAPI.POST("/debug", promptHandler.ReverseProxy(fmt.Sprintf("http://%s:%v/v1/chat/completions", cfg.Proxy.Server.Host, cfg.Proxy.Server.Port)))
		}

		// prompt version api
		versionAPI := authAPI.Group("/version")
		{
			versionAPI.POST("/create", versionHandler.Create)
			versionAPI.GET("/info/:id", versionHandler.GetByID)
			versionAPI.GET("/prompt/:promptId", versionHandler.GetByPromptID)
			versionAPI.GET("/prompt/:promptId/latest", versionHandler.GetLatestByPromptID)
			versionAPI.POST("/update", versionHandler.Update)
			versionAPI.POST("/delete/:id", versionHandler.Delete)
			versionAPI.GET("/list", versionHandler.List)
		}

		// category api
		categoryAPI := authAPI.Group("/category")
		{
			categoryAPI.POST("/create", categoryHandler.Create)
			categoryAPI.GET("/info/:id", categoryHandler.GetByID)
			categoryAPI.GET("/list", categoryHandler.List)
			categoryAPI.POST("/update", categoryHandler.Update)
			categoryAPI.POST("/delete/:id", categoryHandler.Delete)
		}

		// favorites api
		favoritesAPI := authAPI.Group("/favorites")
		{
			favoritesAPI.POST("/add", favoriteHandler.Add)
			favoritesAPI.POST("/remove", favoriteHandler.Remove)
			favoritesAPI.POST("/check", favoriteHandler.Check)
			favoritesAPI.GET("/list", favoriteHandler.List)
		}

		// recently used api
		recentlyUsedAPI := authAPI.Group("/recently-used")
		{
			recentlyUsedAPI.POST("/record", recentlyUsedHandler.Record)
			recentlyUsedAPI.POST("/remove", recentlyUsedHandler.Remove)
			recentlyUsedAPI.GET("/list", recentlyUsedHandler.List)
			recentlyUsedAPI.POST("/clean", recentlyUsedHandler.Clean)
		}
	}

	return r
}
