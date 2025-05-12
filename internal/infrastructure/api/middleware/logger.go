package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()
		errs := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("duration", duration.String()),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		}

		if userID, exists := ctx.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		if errs != "" {
			fields = append(fields, zap.String("errors", errs))
		}

		zap.L().Info("Incoming request", fields...)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("panic recovered",
					zap.Any("error", r),
					zap.String("path", ctx.Request.URL.Path),
					zap.String("method", ctx.Request.Method),
					zap.String("stack", string(debug.Stack())),
				)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()
		ctx.Next()
	}
}
