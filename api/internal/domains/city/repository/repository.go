package repository

import (
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type CityRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *CityRepository {
	return &CityRepository{
		db: store.DB(),
	}
}
