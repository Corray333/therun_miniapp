package repository

import (
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
