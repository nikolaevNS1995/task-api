package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	_ "task-api/docs"
	"task-api/internal/app"
	"task-api/internal/infrastructure/security"
)

// @title Task API
// @version 1.0
// @description API для управления задачами
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	app := fx.New(
		fx.Provide(
			app.NewConfig,
			app.NewLogger,
			app.NewPostgresConnection,
			app.NewTracerProvider,
			app.NewRopositories,
			app.NewUseCases,
			security.NewTokenBlacklist,
			app.NewHandlers,
			gin.New,
		),
		fx.Invoke(
			app.InitTracerProvider,
			app.RegisterRoutes,
			app.RunHTTPServer,
		),
	)
	app.Run()
}
