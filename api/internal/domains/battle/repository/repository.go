package repository

import (
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type BattleRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *BattleRepository {
	return &BattleRepository{
		db: store.DB(),
	}
}
