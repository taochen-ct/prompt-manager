//go:build wireinject
// +build wireinject

package app

import (
	"backend/internal/api/handler"
	"backend/internal/api/middleware"
	"backend/internal/api/router"
	categoryRepo "backend/internal/repository/category"
	favoritesRepo "backend/internal/repository/favorites"
	promptRepo "backend/internal/repository/prompt"
	recentlyUsedRepo "backend/internal/repository/recently_used"
	userRepo "backend/internal/repository/user"
	versionRepo "backend/internal/repository/version"
	categoryService "backend/internal/service/category"
	favoritesService "backend/internal/service/favorites"
	promptService "backend/internal/service/prompt"
	recentlyUsedService "backend/internal/service/recently_used"
	remoteLogService "backend/internal/service/remote_log"
	userService "backend/internal/service/user"
	versionService "backend/internal/service/version"
	"backend/pkg/config"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func wireApp(*config.Config, *lumberjack.Logger, *zap.Logger) (*App, func(), error) {
	panic(
		wire.Build(
			createDB,
			userRepo.CreateRepo,
			userService.CreateUserService,
			handler.CreateUserHandler,
			categoryRepo.CreateCategoryRepo,
			categoryService.CreateCategoryService,
			handler.CreateCategoryHandler,
			favoritesRepo.CreateFavoriteRepo,
			favoritesService.CreateFavoriteService,
			handler.CreateFavoriteHandler,
			recentlyUsedRepo.CreateRecentlyUsedRepo,
			recentlyUsedService.CreateRecentlyUsedService,
			handler.CreateRecentlyUsedHandler,
			promptRepo.CreatePromptRepo,
			versionRepo.CreateVersionRepo,
			promptService.CreatePromptService,
			versionService.CreateVersionService,
			handler.CreatePromptHandler,
			handler.CreatePromptVersionHandler,
			remoteLogService.CreateLogService,
			handler.CreateRemoteLogHandler,
			middleware.CreateRecoveryMiddleware,
			middleware.CreateLoggerMiddleware,
			middleware.CreateCORSMiddleware,
			middleware.CreateJWTMiddleware,
			router.SetupRouter,
			createHttpServer,
			createApp,
		),
	)
}
