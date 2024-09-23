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
	_, err := r.db.NamedExec(`
    UPDATE users 
    SET
        in_app_id = :in_app_id, 
        point_balance = :point_balance, 
        race_balance = :race_balance, 
        key_balance = :key_balance, 
        last_claim = :last_claim, 
        farming_from = :farming_from, 
        max_points = :max_points, 
        farm_time = :farm_time, 
        daily_check_streak = :daily_check_streak, 
        daily_check_last = :daily_check_last 
    WHERE user_id = :user_id
`, user)
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

func (r *UserRepository) CountReferals(userID int64) (refsActivated, refsFrozen, refsClaimed int, err error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1", userID)
	row.Scan(&refsFrozen)

	row = r.db.QueryRow("SELECT refs_claimed FROM users WHERE user_id = $1", userID)
	row.Scan(&refsClaimed)

	return refsActivated, refsFrozen, refsClaimed, nil
}

func (r *UserRepository) ChangeRaceBalance(userID int64, amount int) (int, error) {
	rows, err := r.db.Query("UPDATE users SET race_balance = race_balance + $1 WHERE user_id = $2 RETURNING race_balance", amount, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var balance int
	for rows.Next() {
		err = rows.Scan(&balance)
		if err != nil {
			return 0, err
		}
	}

	return balance, err
}

func (r *UserRepository) ChangePointBalance(userID int64, amount int) (int, error) {
	rows, err := r.db.Query("UPDATE users SET point_balance = point_balance + $1 WHERE user_id = $2 RETURNING point_balance", amount, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var balance int
	for rows.Next() {
		err = rows.Scan(&balance)
		if err != nil {
			return 0, err
		}
	}

	return balance, err
}

func (r *UserRepository) ChangeKeyBalance(userID int64, amount int) (int, error) {
	rows, err := r.db.Query("UPDATE users SET key_balance = key_balance + $1 WHERE user_id = $2 RETURNING key_balance", amount, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var balance int
	for rows.Next() {
		err = rows.Scan(&balance)
		if err != nil {
			return 0, err
		}
	}

	return balance, err
}

func (r *UserRepository) ChangeBalances(userID int64, pointsAmount, raceAmount, keyAmount int) (int, int, int, error) {
	rows, err := r.db.Query("UPDATE users SET point_balance = point_balance + $1, race_balance = race_balance + $2, key_balance = key_balance + $3 WHERE user_id = $4 RETURNING point_balance, race_balance, key_balance", pointsAmount, raceAmount, keyAmount, userID)
	if err != nil {
		return 0, 0, 0, err
	}
	defer rows.Close()

	var pointBalance, raceBalance, keyBalance int
	for rows.Next() {
		err = rows.Scan(&pointBalance, &raceBalance, &keyBalance)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	return pointBalance, raceBalance, keyBalance, nil
}
