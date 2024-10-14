package repository

import (
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type RoundRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *RoundRepository {
	return &RoundRepository{
		db: store.DB(),
	}
}
