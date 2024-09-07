package repository

import (
	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *UserRepository {
	return &UserRepository{
		db: store.DB(),
	}
}

func (r *UserRepository) GetUser(userID int64) (*types.User, error) {
	user := &types.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE user_id = $1", userID)
	return user, err
}

func (r *UserRepository) UpdateUser(user *types.User) error {
	_, err := r.db.NamedExec("UPDATE users SET point_balance = :point_balance, last_claim = :last_claim WHERE user_id = :user_id", user)
	return err
}

func (r *UserRepository) GetRefererID(refCode string) (int64, error) {
	var refererID int64
	err := r.db.Get(&refererID, "SELECT user_id FROM users WHERE ref_code = $1", refCode)
	return refererID, err
}

func (r *UserRepository) CreateUser(user *types.User) error {
	_, err := r.db.NamedExec("INSERT INTO users (user_id, username, avatar, in_app_id, point_balance, race_balance, key_balance, last_claim, max_points, farm_time, ref_code, referer) VALUES (:user_id, :username, :avatar, 0, :point_balance, :race_balance, :key_balance, :last_claim, :max_points, :farm_time, :ref_code, :referer)", user)
	return err
}

// TODO: check if this method works
func (r *UserRepository) CheckIfRefCodeExists(refCode string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE ref_code = $1)", refCode)
	return exists, err
}

func (r *UserRepository) ListReferals(userID int64) ([]types.Referal, error) {
	referals := []types.Referal{}
	err := r.db.Select(&referals, "SELECT avatar, username FROM users WHERE referer = $1", userID)
	return referals, err
}

func (r *UserRepository) CountReferals(userID int64) (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE referer = $1", userID)
	return count, err
}
