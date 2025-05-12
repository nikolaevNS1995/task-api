package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"task-api/pkg/config"
	"time"
)

func RunHTTPServer(lc fx.Lifecycle, engine *gin.Engine, cfg *config.AppConfig, logger *zap.Logger) {
	svr := &http.Server{
		Addr:    cfg.AddressServer,
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("starting HTTP server", zap.String("address", cfg.AddressServer))
				if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server...")
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := svr.Shutdown(shutdownCtx); err != nil {
				logger.Fatal("failed to shutdown server", zap.Error(err))
				return err
			}
			logger.Info("HTTP server shut down gracefully")
			return nil
		},
	})
}
