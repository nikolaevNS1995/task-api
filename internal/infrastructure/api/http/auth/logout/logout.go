package logout

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/pkg/config"
)

func Router(r *gin.Engine, handler *Handler) {
	logoutRouter := r.Group("/api/v1/auth")
	logoutRouter.Use(middleware.AuthMiddleware(handler.cfg, handler.blackList))
	logoutRouter.POST("/logout", handler.Logout)
}

type Handler struct {
	cfg       config.AppConfig
	blackList *security.TokenBlacklist
}

func NewAuthHandler(cfg config.AppConfig, blackList *security.TokenBlacklist) *Handler {
	return &Handler{cfg: cfg, blackList: blackList}
}

// Logout godoc
// @Summary Выход пользователя
// @Description Добавляет access-токен в чёрный список для деактивации. Требует авторизации.
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	createrID, _ := c.Get("user_id")
	authReq := c.GetHeader("Authorization")
	if authReq == "" {
		zap.L().Warn("missing authorization header", zap.Any("creater_id", createrID))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}
	tokenStr, err := security.TokenString(authReq)
	if err != nil {
		zap.L().Warn("invalid token", zap.Any("creater_id", createrID))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	claims, err := security.ParseAccessJWT(h.cfg, tokenStr)
	if err != nil {
		zap.L().Warn("invalid token", zap.Any("creater_id", createrID))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	h.blackList.Add(tokenStr, claims.ExpiresAt.Time)
	zap.L().Info("logged out", zap.String("user_id", claims.UserID.String()), zap.Any("creater_id", createrID))
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
