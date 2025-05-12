package app

import "task-api/internal/usecases"

type UseCases struct {
	taskUseCase    usecases.TaskUseCase
	tagUseCase     usecases.TagUseCase
	commentUseCase usecases.CommentUseCase
	userUseCase    usecases.UserUseCase
	authUseCase    usecases.AuthUseCase
}

func NewUseCases(repos *Repositories) *UseCases {
	return &UseCases{
		taskUseCase:    usecases.NewTasksUseCase(repos.taskRepo),
		tagUseCase:     usecases.NewTagsUseCase(repos.tagRepo),
		commentUseCase: usecases.NewCommentUseCase(repos.commentRepo),
		userUseCase:    usecases.NewUserUseCase(repos.userRepo),
		authUseCase:    usecases.NewAuthUseCase(repos.userRepo, repos.refreshTokenRepo),
	}
}
