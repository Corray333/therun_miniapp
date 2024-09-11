package service

// Errors
var ()

type repository interface {
}

type TaskService struct {
	repo repository
}

func New(repo repository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}
