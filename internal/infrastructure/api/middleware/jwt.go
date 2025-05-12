package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-api/internal/infrastructure/security"
	"task-api/pkg/config"
)

func AuthMiddleware(cfg config.AppConfig, blackListToken *security.TokenBlacklist) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenStr, err := security.TokenString(auth)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if blackListToken.IsBlacklisted(tokenStr) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			return
		}
		claims, err := security.ParseAccessJWT(cfg, tokenStr)
		if err != nil || claims == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Set("user_id", claims.UserID)
		ctx.Next()
	}
}
