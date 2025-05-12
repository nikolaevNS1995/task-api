package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/user"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/internal/usecases"
	"task-api/pkg/config"
)

func Router(r *gin.Engine, handler *Handler, cfg config.AppConfig, blackListToken *security.TokenBlacklist) {
	userRouter := r.Group("api/v1/users")
	userRouter.Use(middleware.AuthMiddleware(cfg, blackListToken))
	{
		userRouter.GET("/:id", handler.GetByID)
		userRouter.GET("/email/:email", handler.GetByEmail)
		userRouter.PUT("/:id", handler.Update)
		userRouter.DELETE("/:id", handler.Delete)
	}
}

type Handler struct {
	useCase usecases.UserUseCase
}

func NewUserHandler(useCase usecases.UserUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// GetByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по его UUID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID пользователя"
// @Success 200 {object} user.UserResponse
// @Router /users/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	createrID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid user id", zap.String("user_id", idStr), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity, err := h.useCase.GetById(c, id)
	if err != nil {
		zap.L().Error("failed get user", zap.String("user_id", id.String()), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success get user", zap.String("user_id", id.String()), zap.Any("creater_id", createrID))
	c.JSON(http.StatusOK, user.FromEntityUser(entity))
}

// GetByEmail godoc
// @Summary Получить пользователя по email
// @Description Возвращает пользователя по email-адресу
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email пользователя"
// @Success 200 {object} user.UserResponse
// @Router /users/email/{email} [get]
func (h *Handler) GetByEmail(c *gin.Context) {
	createrID, _ := c.Get("user_id")
	email := c.Param("email")
	entity, err := h.useCase.GetByEmail(c, email)
	if err != nil {
		zap.L().Error("failed get user", zap.String("email", email), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success get user", zap.String("email", email), zap.Any("creater_id", createrID))
	c.JSON(http.StatusOK, user.FromEntityUser(entity))
}

// Update godoc
// @Summary Обновить пользователя
// @Description Обновляет данные пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID пользователя"
// @Param request body user.UpdateUserRequest true "Данные пользователя"
// @Success 200 {object} user.UserResponse
// @Router /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	createrID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid user id", zap.String("user_id", idStr), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request *user.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid user request", zap.String("user_id", idStr), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := request.ToEntity(id)
	if err := h.useCase.Update(c, entity); err != nil {
		zap.L().Error("failed update user", zap.String("user_id", id.String()), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success update user", zap.String("user_id", id.String()), zap.Any("creater_id", createrID))
	c.JSON(http.StatusOK, user.FromEntityUser(entity))
}

// Delete godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID пользователя"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	createrID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid user id", zap.String("user_id", idStr), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.useCase.Delete(c, id); err != nil {
		zap.L().Error("failed delete user", zap.String("user_id", id.String()), zap.Error(err), zap.Any("creater_id", createrID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success delete user", zap.String("user_id", id.String()), zap.Any("creater_id", createrID))
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
