package service

import (
	"github.com/Corray333/therun_miniapp/internal/domains/task/types"
)

// Errors
var ()

type repository interface {
	GetTasks(userID int64, lang string) ([]*types.Task, error)
}

type external interface {
}

type TaskService struct {
	repo     repository
	external external
}

func New(repo repository, ext external) *TaskService {
	return &TaskService{
		repo:     repo,
		external: ext,
	}
}

func (s *TaskService) GetTasks(userID int64, lang string) ([]*types.Task, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetTasks(userID, lang)
}
