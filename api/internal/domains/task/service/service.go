package service

import (
	"encoding/json"
	"errors"

	"github.com/Corray333/therun_miniapp/internal/domains/task/types"
)

// Errors
var ()

type repository interface {
	GetTasks(userID int64, lang string) ([]*types.Task, error)
	GetTask(taskID int64) (*types.Task, error)
	IsTaskDone(userID, taskID int64) (done bool, err error)

	Claim(userID int64, taskID *types.Task) (done bool, pointsBalance, raceBalance, keysBalance int, err error)
}

type external interface {
	IsUserInChat(chatID, userID int64) (bool, error)
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

func (s *TaskService) CheckTask(userID, taskID int64) (done bool, err error) {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return false, err
	}

	done, err = s.repo.IsTaskDone(userID, taskID)
	if err != nil {
		return false, err
	}
	if done {
		return true, errors.New("task already done")
	}

	switch task.Type {
	case types.TaskTypeTg:
		data := &types.TaskDataTG{}
		if err := json.Unmarshal(task.Data, data); err != nil {
			return false, err
		}
		done, err := s.external.IsUserInChat(data.ChatID, userID)
		if err != nil {
			return false, err
		}
		return done, nil
	case types.TaskTypeUncheckable:
		return true, nil
	default:
		return false, nil
	}
}

func (s *TaskService) Claim(userID, taskID int64) (done bool, pointsBalance, raceBalance, keysBalance int, err error) {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return false, 0, 0, 0, err
	}

	done, err = s.repo.IsTaskDone(userID, taskID)
	if err != nil {
		return false, 0, 0, 0, err
	}
	if done {
		return true, 0, 0, 0, errors.New("task already done")
	}

	switch task.Type {
	case types.TaskTypeTg:
		data := &types.TaskDataTG{}
		if err := json.Unmarshal(task.Data, data); err != nil {
			return false, 0, 0, 0, err
		}

		done, err = s.external.IsUserInChat(data.ChatID, userID)
		if err != nil {
			return false, 0, 0, 0, err
		}
	case types.TaskTypeUncheckable:
		done = true

	}

	if done {
		done, pointsBalance, raceBalance, keysBalance, err = s.repo.Claim(userID, task)
		if err != nil {
			return false, 0, 0, 0, err
		}
	}

	return done, pointsBalance, raceBalance, keysBalance, nil
}
