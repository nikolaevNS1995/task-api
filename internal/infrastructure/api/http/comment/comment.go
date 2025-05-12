package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/comment"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/internal/usecases"
	"task-api/pkg/config"
)

func Router(r *gin.Engine, handler *Handler, cfg config.AppConfig, blackListToken *security.TokenBlacklist) {
	commentRouter := r.Group("/api/v1/comments")
	commentRouter.Use(middleware.AuthMiddleware(cfg, blackListToken))
	{
		commentRouter.GET("/", handler.GetAll)
		commentRouter.POST("/", handler.Create)
		commentRouter.GET("/:id", handler.GetById)
		commentRouter.PUT("/:id", handler.Update)
		commentRouter.DELETE("/:id", handler.Delete)
	}
}

type Handler struct {
	useCase usecases.CommentUseCase
}

func NewCommentHandler(useCase usecases.CommentUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// GetAll godoc
// @Summary Получить все комментарии
// @Description Возвращает все комментарии, доступные пользователю
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} comment.CommentResponse
// @Router /comments [get]
func (h *Handler) GetAll(c *gin.Context) {
	userID, _ := c.Get("user_id")
	comments, err := h.useCase.GetAll(c)
	if err != nil {
		zap.L().Error("failed get comments", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var output []*comment.CommentResponse
	for _, m := range comments {
		output = append(output, comment.FromModelComment(m))
	}
	zap.L().Info("success get comments", zap.Int("count_comments", len(comments)), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, output)
}

// Create godoc
// @Summary Создать комментарий
// @Description Создаёт новый комментарий к задаче
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body comment.CreateCommentRequest true "Данные комментария"
// @Success 200 {object} comment.CommentResponse
// @Router /comments [post]
func (h *Handler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var request comment.CreateCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid comment request", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := request.ToEntity()

	create, err := h.useCase.Create(c, entity)
	if err != nil {
		zap.L().Error("failed create comment", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success create comment", zap.String("comment_id", create.Comment.ID.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, comment.FromModelComment(create))
}

// GetById godoc
// @Summary Получить комментарий по ID
// @Description Возвращает комментарий по его ID
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID комментария"
// @Success 200 {object} comment.CommentResponse
// @Router /comments/{id} [get]
func (h *Handler) GetById(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid comment id", zap.String("comment_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model, err := h.useCase.GetByID(c, id)
	if err != nil {
		zap.L().Error("failed get comment", zap.String("comment_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comment.FromModelComment(model))
}

// Update godoc
// @Summary Обновить комментарий
// @Description Обновляет содержимое комментария по его ID
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID комментария"
// @Param request body comment.UpdateCommentRequest true "Новые данные комментария"
// @Success 200 {object} comment.CommentResponse
// @Router /comments/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid comment id", zap.String("comment_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request *comment.UpdateCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid comment request", zap.String("comment_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := request.ToEntity(id)
	model, err := h.useCase.Update(c, entity)
	if err != nil {
		zap.L().Error("failed update comment", zap.String("comment_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success update comment", zap.String("comment_id", id.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, comment.FromModelComment(model))
}

// Delete godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий по ID
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID комментария"
// @Success 200 {object} map[string]string
// @Router /comments/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid comment id", zap.String("comment_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.useCase.Delete(c, id); err != nil {
		zap.L().Error("failed delete comment", zap.String("comment_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("success delete comment", zap.String("comment_id", id.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
