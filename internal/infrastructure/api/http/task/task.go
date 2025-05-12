package task

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"task-api/internal/adapters/api/task"
	"task-api/internal/domain/entities"
	"task-api/internal/infrastructure/api/middleware"
	"task-api/internal/infrastructure/security"
	"task-api/internal/usecases"
	"task-api/pkg/config"
)

func Router(router *gin.Engine, handler *Handler, cfg config.AppConfig, blackListToken *security.TokenBlacklist) {
	taskRouter := router.Group("/api/v1/tasks")
	taskRouter.Use(middleware.AuthMiddleware(cfg, blackListToken))
	{
		taskRouter.GET("", handler.GetTasks)
		taskRouter.GET("/:id", handler.GetTask)
		taskRouter.POST("", handler.CreateTask)
		taskRouter.PUT("/:id", handler.UpdateTask)
		taskRouter.DELETE("/:id", handler.DeleteTask)
		taskRouter.POST("/:id/tags", handler.AddTags)
		taskRouter.DELETE("/:id/tags", handler.DeleteTags)

	}
}

type Handler struct {
	useCase usecases.TaskUseCase
}

func NewTaskHandler(useCase usecases.TaskUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// CreateTask godoc
// @Summary Создать задачу
// @Description Создание новой задачи пользователем
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body task.CreateTaskRequest true "Данные задачи"
// @Success 201 {object} task.TaskResponse
// @Router /tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	var request task.CreateTaskRequest
	userID, exists := c.Get("user_id")
	if !exists {
		zap.L().Warn("no user id provided", zap.Error(errors.New("no user id provided")))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user id provided"})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid request", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := request.ToEntity(userID.(uuid.UUID))
	model, err := h.useCase.Create(c, entity)
	if err != nil {
		zap.L().Error("failed to create task", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("task created", zap.String("task_id", model.Task.ID.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusCreated, task.FromModelTask(model))
}

// GetTask godoc
// @Summary Получить задачу
// @Description Получение задачи по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID задачи"
// @Success 200 {object} task.TaskResponse
// @Router /tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid request", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model, err := h.useCase.GetTask(c, id)
	if err != nil {
		zap.L().Error("failed to get task", zap.String("task_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("task get", zap.String("task_id", model.Task.ID.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, task.FromModelTask(model))
}

// GetTasks godoc
// @Summary Получить список задач
// @Description Получение всех задач текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} task.TaskAllResponse
// @Router /tasks [get]
func (h *Handler) GetTasks(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tasks, err := h.useCase.GetTasksByUserID(c, userID.(uuid.UUID))
	if err != nil {
		zap.L().Error("failed to get tasks", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var output []*task.TaskAllResponse
	for _, model := range tasks {
		output = append(output, task.FromModelTaskForAll(model))
	}
	zap.L().Info("tasks get", zap.Int("count", len(tasks)), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, output)
}

// UpdateTask godoc
// @Summary Обновить задачу
// @Description Обновление существующей задачи по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID задачи"
// @Param request body task.UpdateTaskRequest true "Новые данные задачи"
// @Success 200 {object} task.TaskResponse
// @Router /tasks/{id} [put]
func (h *Handler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid updated task ID", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request task.UpdateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid updated task request", zap.String("task_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := request.ToEntity(id)
	model, err := h.useCase.Update(c, entity)
	if err != nil {
		zap.L().Error("failed to update task", zap.String("task_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("task updated", zap.String("task_id", model.Task.ID.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, task.FromModelTask(model))
}

// DeleteTask godoc
// @Summary Удалить задачу
// @Description Удаление задачи по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID задачи"
// @Success 200 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid deleted task ID", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.useCase.Delete(c, id); err != nil {
		zap.L().Error("failed to delete task", zap.String("task_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("task deleted", zap.String("task_id", id.String()), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

// AddTags godoc
// @Summary Добавить теги к задаче
// @Description Привязка тегов к задаче
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID задачи"
// @Param request body []task.TagRequest true "Список тегов"
// @Success 200 {object} map[string]string
// @Router /tasks/{id}/tags [post]
func (h *Handler) AddTags(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid task ID for add tags", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request []*task.TagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid request for add tags", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tags []*entities.Tag
	for _, tag := range request {
		entity := tag.ToEntity()
		tags = append(tags, entity)
	}

	if err := h.useCase.AddTags(c, id, tags); err != nil {
		zap.L().Error("failed to add tags", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("tags added", zap.String("task_id", id.String()), zap.Int("tag_count", len(tags)), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "tags added"})
}

// DeleteTags godoc
// @Summary Удалить теги у задачи
// @Description Удаление тегов из задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID задачи"
// @Param request body []task.TagRequest true "Список тегов для удаления"
// @Success 200 {object} map[string]string
// @Router /tasks/{id}/tags [delete]
func (h *Handler) DeleteTags(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.L().Warn("invalid task ID for delete tags", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request []*task.TagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Warn("invalid request for delete tags", zap.String("task_id", idStr), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tags []*entities.Tag
	for _, tag := range request {
		entity := tag.ToEntity()
		tags = append(tags, entity)
	}

	if err := h.useCase.RemoveTags(c, id, tags); err != nil {
		zap.L().Error("failed to remove tags", zap.String("task_id", id.String()), zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zap.L().Info("tags removed", zap.String("task_id", id.String()), zap.Int("tag_count", len(tags)), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "tags removed"})
}
