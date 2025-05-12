package app

import (
	"task-api/internal/infrastructure/api/http/auth/login"
	"task-api/internal/infrastructure/api/http/auth/logout"
	"task-api/internal/infrastructure/api/http/auth/me"
	"task-api/internal/infrastructure/api/http/auth/refresh"
	"task-api/internal/infrastructure/api/http/auth/registr"
	"task-api/internal/infrastructure/api/http/comment"
	"task-api/internal/infrastructure/api/http/tag"
	"task-api/internal/infrastructure/api/http/task"
	"task-api/internal/infrastructure/api/http/user"
	"task-api/internal/infrastructure/security"
	"task-api/pkg/config"
)

type Handlers struct {
	taskHandler    *task.Handler
	tagHandler     *tag.Handler
	commentHandler *comment.Handler
	userHandler    *user.Handler
	loginHandler   *login.Handler
	registHandler  *registr.Handler
	logoutHandler  *logout.Handler
	meHandler      *me.Handler
	refreshHandler *refresh.Handler
}

func NewHandlers(useCase *UseCases, cfg *config.AppConfig, blackListToken *security.TokenBlacklist) *Handlers {
	return &Handlers{
		taskHandler:    task.NewTaskHandler(useCase.taskUseCase),
		tagHandler:     tag.NewTagHandler(useCase.tagUseCase),
		commentHandler: comment.NewCommentHandler(useCase.commentUseCase),
		userHandler:    user.NewUserHandler(useCase.userUseCase),
		//authHandler
		loginHandler:   login.NewAuthHandler(useCase.authUseCase, *cfg),
		registHandler:  registr.NewAuthHandler(useCase.authUseCase),
		logoutHandler:  logout.NewAuthHandler(*cfg, blackListToken),
		meHandler:      me.NewAuthHandler(useCase.userUseCase),
		refreshHandler: refresh.NewAuthHandler(useCase.authUseCase, *cfg),
	}
}
