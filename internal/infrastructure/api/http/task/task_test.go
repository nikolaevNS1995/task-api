package task_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"task-api/internal/adapters/api/task"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	handler "task-api/internal/infrastructure/api/http/task"
	"task-api/internal/usecases/mocks"
	"testing"
)

func TestHandler_GetTask_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)

	h := handler.NewTaskHandler(mockUseCase)

	taskID := uuid.New()
	userId := uuid.New()

	expectedModel := &models.Task{
		Task: entities.Task{
			ID:    taskID,
			Title: "New Task",
		},
	}

	mockUseCase.EXPECT().GetTask(gomock.Any(), taskID).Return(expectedModel, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user_id", userId)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
	c.Request = req

	c.Params = gin.Params{{Key: "id", Value: taskID.String()}}

	h.GetTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":"`+taskID.String()+`"`)
}

func TestHandler_GetTask_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)
	h := handler.NewTaskHandler(mockUseCase)

	taskID := uuid.New()
	userId := uuid.New()

	mockUseCase.EXPECT().GetTask(gomock.Any(), taskID).Return(nil, errors.New("task not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userId)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: taskID.String()}}

	h.GetTask(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "task not found")
}

func TestHandler_CreateTask_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)
	h := handler.NewTaskHandler(mockUseCase)

	userID := uuid.New()
	taskID := uuid.New()

	input := task.CreateTaskRequest{
		Title:       "Test task",
		Description: "Some description",
	}
	expected := &models.Task{
		Task: entities.Task{
			ID:          taskID,
			Title:       input.Title,
			Description: input.Description,
			CreatedBy:   userID,
		},
	}

	mockUseCase.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&entities.Task{})).
		DoAndReturn(func(_ context.Context, actual *entities.Task) (*models.Task, error) {
			expectedEntity := input.ToEntity(userID)
			require.Equal(t, expectedEntity.Title, actual.Title)
			require.Equal(t, expectedEntity.Description, actual.Description)
			require.Equal(t, expectedEntity.CreatedBy, actual.CreatedBy)
			return expected, nil
		})

	// Сериализуем input
	body, err := json.Marshal(input)
	require.NoError(t, err)

	// Создаём контекст Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызываем handler
	h.CreateTask(c)

	// Проверяем ответ
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), taskID.String())
}

func TestHandler_CreateTask_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)
	h := handler.NewTaskHandler(mockUseCase)

	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer([]byte(`invalid json`)))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestHandler_CreateTask_UseCaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)
	h := handler.NewTaskHandler(mockUseCase)

	userID := uuid.New()

	input := task.CreateTaskRequest{
		Title:       "Task",
		Description: "Some description",
	}

	mockUseCase.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&entities.Task{})).
		Return(nil, errors.New("create failed"))

	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateTask(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "create failed")
}

func TestHandler_CreateTask_MissingUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)
	h := handler.NewTaskHandler(mockUseCase)

	input := task.CreateTaskRequest{
		Title:       "Task",
		Description: "Some description",
	}

	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Не устанавливаем user_id: c.Set("user_id", userID)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateTask(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "no user id provided")
}

func TestHandler_UpdateTask_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)
	h := handler.NewTaskHandler(mockUseCase)

	userID := uuid.New()
	taskID := uuid.New()

	input := task.UpdateTaskRequest{
		Title:       "Updated title",
		Description: "Updated description",
		Status:      "in_progress",
	}

	updatedTask := &models.Task{
		Task: entities.Task{
			ID:          taskID,
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			CreatedBy:   userID,
		},
	}

	mockUseCase.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, task *entities.Task) (*models.Task, error) {
			expectedEntity := input.ToEntity(userID)
			assert.Equal(t, expectedEntity.CreatedBy, task.CreatedBy)
			assert.Equal(t, expectedEntity.Title, task.Title)
			assert.Equal(t, expectedEntity.Description, task.Description)
			assert.Equal(t, expectedEntity.Status, task.Status)
			return updatedTask, nil
		})

	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Устанавливаем param
	c.Params = gin.Params{gin.Param{Key: "id", Value: taskID.String()}}

	h.UpdateTask(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp entities.Task
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, updatedTask.Task.ID, resp.ID)
	assert.Equal(t, updatedTask.Task.Title, resp.Title)
	assert.Equal(t, updatedTask.Task.Description, resp.Description)
	assert.Equal(t, updatedTask.Task.Status, resp.Status)
}

func TestHandler_GetTasks_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTaskUseCase(ctrl)

	h := handler.NewTaskHandler(mockUseCase)

	userId := uuid.New()
	taskID1 := uuid.New()
	taskID2 := uuid.New()

	expectedModel := []*models.TasksWishTags{
		{
			Task: entities.Task{
				ID:          taskID1,
				Title:       "New Task1",
				Description: "New Description1",
				Status:      "in_progress",
				CreatedBy:   userId,
			},
		},
		{
			Task: entities.Task{
				ID:          taskID2,
				Title:       "New Task2",
				Description: "New Description2",
				Status:      "in_progress",
				CreatedBy:   userId,
			},
		},
	}

	mockUseCase.EXPECT().
		GetTasksByUserID(gomock.Any(), userId).
		Return(expectedModel, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user_id", userId)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks/", nil)
	c.Request = req

	h.GetTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse []*task.TaskAllResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expected := []*task.TaskAllResponse{
		{
			ID:          taskID1,
			Title:       "New Task1",
			Description: "New Description1",
			Status:      "in_progress",
			Tags:        nil, // или []string{}, если нужно
		},
		{
			ID:          taskID2,
			Title:       "New Task2",
			Description: "New Description2",
			Status:      "in_progress",
			Tags:        nil,
		},
	}

	assert.Equal(t, expected, actualResponse)
}
