package tag

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/tag"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/internal/usecases"
	"task-api/pkg/config"
)

func Router(router *gin.Engine, handler *Handler, cfg config.AppConfig, blackListToken *security.TokenBlacklist) {
	tagRouter := router.Group("/api/v1/tags")
	tagRouter.Use(middleware.AuthMiddleware(cfg, blackListToken))
	{
		tagRouter.GET("/", handler.GetTags)
		tagRouter.POST("/", handler.Create)
		tagRouter.GET("/:id", handler.GetTag)
		tagRouter.PUT("/:id", handler.Update)
		tagRouter.DELETE("/:id", handler.Delete)
	}
}

type Handler struct {
	useCase usecases.TagUseCase
}

func NewTagHandler(useCase usecases.TagUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// GetTags godoc
// @Summary Получить все теги
// @Description Получает список всех тегов пользователя
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} tag.TagResponse
// @Router /tags [get]
func (h *Handler) GetTags(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tags, err := h.useCase.GetTags(c)
	if err != nil {
		zap.L().Error("failed get tags", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var output []*tag.TagResponse
	for _, entity := range tags {
		output = append(output, tag.FromEntityTag(entity))
	}
	zap.L().Info("success get tags", zap.Int("count_tags", len(tags)), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, output)

}

// Create godoc
// @Summary Создать тег
// @Description Создаёт новый тег
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tag.CreateTagRequest true "Данные нового тега"
// @Success 200 {object} tag.TagResponse
// @Router /tags [post]
func (h *Handler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var request tag.CreateTagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid tag request", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity, err := h.useCase.Create(c, request.ToEntity())
	if err != nil {
		zap.L().Error("failed create tag", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success create tag", zap.String("tag_id", entity.ID.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, tag.FromEntityTag(entity))
}

// GetTag godoc
// @Summary Получить тег
// @Description Получает тег по ID
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID тега"
// @Success 200 {object} tag.TagResponse
// @Router /tags/{id} [get]
func (h *Handler) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid tag id", zap.String("tag_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity, err := h.useCase.GetTag(c, id)
	if err != nil {
		zap.L().Error("failed get tag", zap.String("tag_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success get tag", zap.String("tag_id", id.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, tag.FromEntityTag(entity))
}

// Update godoc
// @Summary Обновить тег
// @Description Обновляет тег по ID
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID тега"
// @Param request body tag.UpdateTagRequest true "Новые данные тега"
// @Success 200 {object} tag.TagResponse
// @Router /tags/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid tag id", zap.String("tag_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request tag.UpdateTagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid tag request", zap.String("tag_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := request.ToEntity(id)
	if err := h.useCase.Update(c, entity); err != nil {
		zap.L().Error("failed update tag", zap.String("tag_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success update tag", zap.String("tag_id", id.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, tag.FromEntityTag(entity))
}

// Delete godoc
// @Summary Удалить тег
// @Description Удаляет тег по ID
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID тега"
// @Success 200 {object} map[string]string
// @Router /tags/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid tag id", zap.String("tag_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.useCase.Delete(c, id); err != nil {
		zap.L().Error("failed delete tag", zap.String("tag_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success delete tag", zap.String("tag_id", id.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
}
