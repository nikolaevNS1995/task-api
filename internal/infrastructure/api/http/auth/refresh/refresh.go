package refresh

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/auth"
	"task-api/internal/infrastructure/security"
	"task-api/internal/usecases"
	"task-api/pkg/config"
	"time"
)

func Router(r *gin.Engine, handler *Handler) {
	refreshRouter := r.Group("/api/v1/auth")
	refreshRouter.POST("/refresh", handler.Refresh)
}

type Handler struct {
	useCase usecases.AuthUseCase
	cfg     config.AppConfig
}

func NewAuthHandler(useCase usecases.AuthUseCase, cfg config.AppConfig) *Handler {
	return &Handler{useCase: useCase, cfg: cfg}
}

// Refresh godoc
// @Summary Refresh access token
// @Description Обновляет access token используя refresh token
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param refreshToken body auth.RefreshRequest true "Refresh Token"
// @Success 200 {object} auth.LoginResponse
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var request auth.RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid refresh request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err := h.useCase.GetRefreshToken(c, request.RefreshToken)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) {
		zap.L().Warn("invalid refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if err := h.useCase.DeleteRefreshToken(c, request.RefreshToken); err != nil {
		zap.L().Warn("failed delete refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete refresh token"})
		return
	}

	signedToken, err := security.CreateAccessJWT(h.cfg, refreshToken.UserID)
	if err != nil {
		zap.L().Warn("failed create access token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expiresAt := time.Now().Add(h.cfg.Auth.JWTRefreshExpiry)

	newRefreshToken, err := h.useCase.CreateRefreshToken(c, refreshToken.UserID, expiresAt)
	if err != nil {
		zap.L().Warn("failed create refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}
	zap.L().Info("success refresh token", zap.String("user_id", refreshToken.UserID.String()))
	c.JSON(http.StatusOK, auth.LoginResponse{AccessToken: signedToken, RefreshToken: newRefreshToken.Token.String()})
}
