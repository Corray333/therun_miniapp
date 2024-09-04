package repository

import (
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type FarmingRepository struct {
	db *sqlx.DB
}

func New(store storage.Storage) *FarmingRepository {
	return &FarmingRepository{
		db: store.DB(),
	}
}
