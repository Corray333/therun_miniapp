package repository

import (
	"github.com/Corray333/therun_miniapp/internal/domains/task/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *TaskRepository {
	return &TaskRepository{
		db: store.DB(),
	}
}

func (r *TaskRepository) GetTasks(userID int64, lang string) ([]*types.Task, error) {

	tasks := []*types.Task{}
	err := r.db.Select(&tasks, `
	SELECT t.*, tt.description
	FROM tasks t
	LEFT JOIN task_translate tt ON t.task_id = tt.task_id AND tt.lang = $2
	WHERE t.expire_at > EXTRACT(EPOCH FROM NOW())
	AND NOT EXISTS (
		SELECT 1
		FROM user_task_done ut
		WHERE ut.task_id = t.task_id
		AND ut.user_id = $1
	) 
	`, userID, lang)
	return tasks, err
}
