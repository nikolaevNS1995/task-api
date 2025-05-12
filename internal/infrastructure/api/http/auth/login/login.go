package login

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
	loginRouter := r.Group("/api/v1/auth")
	loginRouter.POST("/login", handler.Login)
}

type Handler struct {
	useCase usecases.AuthUseCase
	cfg     config.AppConfig
}

func NewAuthHandler(useCase usecases.AuthUseCase, cfg config.AppConfig) *Handler {
	return &Handler{useCase: useCase, cfg: cfg}
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Проверка email и пароля. Возвращает access и refresh токены.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Данные пользователя для входа"
// @Success 200 {object} auth.LoginResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var request auth.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.useCase.Login(c, request.Email, request.Password)
	if err != nil {
		zap.L().Warn("failed login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	signedToken, err := security.CreateAccessJWT(h.cfg, user.ID)
	if err != nil {
		zap.L().Warn("failed create access token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expiresAt := time.Now().Add(h.cfg.Auth.JWTRefreshExpiry)

	refreshToken, err := h.useCase.CreateRefreshToken(c, user.ID, expiresAt)
	if err != nil {
		zap.L().Warn("failed create refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}
	zap.L().Info("success login", zap.String("user_id", user.ID.String()))
	c.JSON(http.StatusOK, auth.LoginResponse{AccessToken: signedToken, RefreshToken: refreshToken.Token.String()})
}
