package me

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/auth"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/internal/usecases"
	"task-api/pkg/config"
)

func Router(r *gin.Engine, handler *Handler, cfg config.AppConfig, blackListToken *security.TokenBlacklist) {
	meRouter := r.Group("api/v1/auth")
	meRouter.Use(middleware.AuthMiddleware(cfg, blackListToken))
	meRouter.POST("/me", handler.Me)
}

type Handler struct {
	useCase usecases.UserUseCase
}

func NewAuthHandler(useCase usecases.UserUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Me godoc
// @Summary Get current user
// @Description Возвращает данные текущего пользователя
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} auth.MeResponse
// @Failure 401 {object} map[string]string
// @Router /auth/me [post]
func (h *Handler) Me(c *gin.Context) {
	userIdRaw, exists := c.Get("user_id")
	if !exists {
		zap.L().Warn("missing user id", zap.Any("creater_id", userIdRaw))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		zap.L().Warn("invalid user id", zap.Any("creater_id", userIdRaw))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	usr, err := h.useCase.GetById(c, userId)
	if err != nil {
		zap.L().Warn("failed get user", zap.String("user_id", userId.String()), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	zap.L().Info("success get user", zap.String("user_id", userId.String()))
	c.JSON(http.StatusOK, auth.MeResponse{
		ID:        usr.ID,
		Name:      usr.Name,
		Email:     usr.Email,
		CreatedAt: usr.CreatedAt,
	})
}
