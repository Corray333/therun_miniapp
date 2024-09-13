package repository

import (
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type FarmingRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *FarmingRepository {
	return &FarmingRepository{
		db: store.DB(),
	}
}

// StartFarming(userID int64, startTime int64) error
// Claim(userID int64, pointBalance int, lastClaim int64) error

func (r *FarmingRepository) StartFarming(userID int64, startTime int64) error {
	_, err := r.db.Exec("UPDATE users SET farming_from = $1 WHERE user_id = $2", startTime, userID)
	return err
}

func (r *FarmingRepository) Claim(userID int64, pointBalance int, lastClaim int64) error {
	_, err := r.db.Exec("UPDATE users SET point_balance = $1, last_claim = $2 WHERE user_id = $3", pointBalance, lastClaim, userID)
	return err
}

// func (r *FarmingRepository) FindAvailibleClaim(){}
