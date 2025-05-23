// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecases/task.go
//
// Generated by this command:
//
//	mockgen -source=internal/usecases/task.go -destination=internal/usecases/mocks/task_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	models "task-api/internal/adapters/models"
	entities "task-api/internal/domain/entities"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockTaskUseCase is a mock of TaskUseCase interface.
type MockTaskUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockTaskUseCaseMockRecorder
	isgomock struct{}
}

// MockTaskUseCaseMockRecorder is the mock recorder for MockTaskUseCase.
type MockTaskUseCaseMockRecorder struct {
	mock *MockTaskUseCase
}

// NewMockTaskUseCase creates a new mock instance.
func NewMockTaskUseCase(ctrl *gomock.Controller) *MockTaskUseCase {
	mock := &MockTaskUseCase{ctrl: ctrl}
	mock.recorder = &MockTaskUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskUseCase) EXPECT() *MockTaskUseCaseMockRecorder {
	return m.recorder
}

// AddTags mocks base method.
func (m *MockTaskUseCase) AddTags(ctx context.Context, taskID uuid.UUID, tags []*entities.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTags", ctx, taskID, tags)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTags indicates an expected call of AddTags.
func (mr *MockTaskUseCaseMockRecorder) AddTags(ctx, taskID, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTags", reflect.TypeOf((*MockTaskUseCase)(nil).AddTags), ctx, taskID, tags)
}

// Create mocks base method.
func (m *MockTaskUseCase) Create(ctx context.Context, task *entities.Task) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, task)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTaskUseCaseMockRecorder) Create(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTaskUseCase)(nil).Create), ctx, task)
}

// Delete mocks base method.
func (m *MockTaskUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTaskUseCaseMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTaskUseCase)(nil).Delete), ctx, id)
}

// GetTask mocks base method.
func (m *MockTaskUseCase) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", ctx, id)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskUseCaseMockRecorder) GetTask(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskUseCase)(nil).GetTask), ctx, id)
}

// GetTasks mocks base method.
func (m *MockTaskUseCase) GetTasks(ctx context.Context) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasks", ctx)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasks indicates an expected call of GetTasks.
func (mr *MockTaskUseCaseMockRecorder) GetTasks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasks", reflect.TypeOf((*MockTaskUseCase)(nil).GetTasks), ctx)
}

// GetTasksByUserID mocks base method.
func (m *MockTaskUseCase) GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TasksWishTags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksByUserID", ctx, userID)
	ret0, _ := ret[0].([]*models.TasksWishTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksByUserID indicates an expected call of GetTasksByUserID.
func (mr *MockTaskUseCaseMockRecorder) GetTasksByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksByUserID", reflect.TypeOf((*MockTaskUseCase)(nil).GetTasksByUserID), ctx, userID)
}

// RemoveTags mocks base method.
func (m *MockTaskUseCase) RemoveTags(ctx context.Context, taskID uuid.UUID, tags []*entities.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveTags", ctx, taskID, tags)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveTags indicates an expected call of RemoveTags.
func (mr *MockTaskUseCaseMockRecorder) RemoveTags(ctx, taskID, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveTags", reflect.TypeOf((*MockTaskUseCase)(nil).RemoveTags), ctx, taskID, tags)
}

// Update mocks base method.
func (m *MockTaskUseCase) Update(ctx context.Context, task *entities.Task) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, task)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockTaskUseCaseMockRecorder) Update(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTaskUseCase)(nil).Update), ctx, task)
}
