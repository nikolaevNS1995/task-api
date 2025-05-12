package registr

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/user"
	"task-api/internal/usecases"
)

func Router(r *gin.Engine, handler *Handler) {
	registrRouter := r.Group("/api/v1/auth")
	registrRouter.POST("/registration", handler.Regist)
}

type Handler struct {
	useCase usecases.AuthUseCase
}

func NewAuthHandler(useCase usecases.AuthUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Regist Register godoc
// @Summary User Registration
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "New user data"
// @Success 200 {object} user.UserResponse
// @Router /auth/registration [post]
func (h *Handler) Regist(c *gin.Context) {
	var request *user.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid user request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity, err := h.useCase.Register(c, request.ToEntity())
	if err != nil {
		zap.L().Warn("failed create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success create user", zap.String("user_id", entity.ID.String()))
	c.JSON(http.StatusOK, user.FromEntityUser(entity))
}
