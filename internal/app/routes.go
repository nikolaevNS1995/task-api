package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"task-api/internal/infrastructure/api/http/auth/login"
	"task-api/internal/infrastructure/api/http/auth/logout"
	"task-api/internal/infrastructure/api/http/auth/me"
	"task-api/internal/infrastructure/api/http/auth/refresh"
	"task-api/internal/infrastructure/api/http/auth/registr"
	"task-api/internal/infrastructure/api/http/comment"
	"task-api/internal/infrastructure/api/http/tag"
	"task-api/internal/infrastructure/api/http/task"
	"task-api/internal/infrastructure/api/http/user"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/pkg/config"
)

func RegisterRoutes(router *gin.Engine, cfg *config.AppConfig, handers *Handlers, blackListToken *security.TokenBlacklist) {
	// Middleware
	router.Use(middleware.TracingMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	task.Router(router, handers.taskHandler, *cfg, blackListToken)
	tag.Router(router, handers.tagHandler, *cfg, blackListToken)
	comment.Router(router, handers.commentHandler, *cfg, blackListToken)
	user.Router(router, handers.userHandler, *cfg, blackListToken)
	//Auth Routes
	login.Router(router, handers.loginHandler)
	registr.Router(router, handers.registHandler)
	logout.Router(router, handers.logoutHandler)
	me.Router(router, handers.meHandler, *cfg, blackListToken)
	refresh.Router(router, handers.refreshHandler)
}
