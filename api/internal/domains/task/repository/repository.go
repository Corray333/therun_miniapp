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

func (r *TaskRepository) GetTask(taskID int64) (*types.Task, error) {
	task := &types.Task{}
	err := r.db.Get(task, `
	SELECT t.*, tt.description
	FROM tasks t
	LEFT JOIN task_translate tt ON t.task_id = tt.task_id AND tt.lang = 'en'
	WHERE t.task_id = $1
	`, taskID)
	return task, err
}

func (r *TaskRepository) IsTaskDone(userID, taskID int64) (done bool, err error) {
	var count int
	err = r.db.Get(&count, "SELECT COUNT(*) FROM user_task_done WHERE user_id = $1 AND task_id = $2", userID, taskID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TaskRepository) Claim(userID int64, task *types.Task) (done bool, pointBalance, raceBalance, keyBalance int, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return false, 0, 0, 0, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("UPDATE users SET point_balance = point_balance + $1, race_balance = race_balance + $2, red_key_balance = red_key_balance + $3 WHERE user_id = $4 RETURNING point_balance, race_balance, red_key_balance", task.PointsReward, task.RaceReward, task.KeysReward, userID)
	if err != nil {
		return false, 0, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pointBalance, &raceBalance, &keyBalance)
		if err != nil {
			return false, 0, 0, 0, err
		}
	}

	_, err = tx.Exec("INSERT INTO user_task_done (user_id, task_id) VALUES ($1, $2)", userID, task.ID)
	if err != nil {
		return false, 0, 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return false, 0, 0, 0, err
	}
	return true, pointBalance, raceBalance, keyBalance, nil
}
