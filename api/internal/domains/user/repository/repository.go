package repository

import (
	"context"
	"fmt"

	"github.com/Corray333/therun_miniapp/internal/domains/user/service"
	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	global_types "github.com/Corray333/therun_miniapp/internal/types"
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

func (r *UserRepository) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, global_types.TxKey{}, tx), nil
}

func (r *UserRepository) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(global_types.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx.Commit()
}

func (r *UserRepository) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value(global_types.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx.Rollback()
}

func (r *UserRepository) GetTx(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(global_types.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
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
	_, err := r.db.NamedExec("INSERT INTO users (user_id, username, avatar, in_app_id, point_balance, race_balance, key_balance, last_claim, max_points, farm_time, ref_code, referer, is_premium, is_activated) VALUES (:user_id, :username, :avatar, 0, :point_balance, :race_balance, :key_balance, :last_claim, :max_points, :farm_time, :ref_code, :referer, :is_premium, :is_activated)", user)
	return err
}

// TODO: check if this method works
func (r *UserRepository) CheckIfRefCodeExists(refCode string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE ref_code = $1)", refCode)
	return exists, err
}

func (r *UserRepository) ListActivatedReferals(userID int64) ([]types.Referal, error) {
	referals := []types.Referal{}
	err := r.db.Select(&referals, "SELECT avatar, username FROM users WHERE referer = $1 AND is_activated = true", userID)
	return referals, err
}

func (r *UserRepository) ListNotActivatedReferals(userID int64) ([]types.Referal, error) {
	referals := []types.Referal{}
	err := r.db.Select(&referals, "SELECT avatar, username FROM users WHERE referer = $1 AND is_activated = false", userID)
	return referals, err
}

func (r *UserRepository) CountReferals(userID int64) (refsActivated, refsFrozen, refsPremiumActivatedNotClaimed, refsPremiumFrozenNotClaimed, refsActivatedNotClaimed, refsFrozenNotClaimed int, err error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = true", userID)
	if err := row.Scan(&refsActivated); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = false", userID)
	if err := row.Scan(&refsFrozen); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsPremiumActivatedNotClaimed); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = false AND ref_claimed = false", userID)
	if err := row.Scan(&refsPremiumFrozenNotClaimed); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = false AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsActivatedNotClaimed); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = false AND is_activated = false AND ref_claimed = false", userID)
	if err := row.Scan(&refsFrozenNotClaimed); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	return refsActivated, refsFrozen, refsPremiumActivatedNotClaimed, refsPremiumFrozenNotClaimed, refsActivatedNotClaimed, refsFrozenNotClaimed, nil
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

func (r *UserRepository) SetPremium(userID int64, isPremium bool) error {
	_, err := r.db.Exec("UPDATE users SET is_premium = $1 WHERE user_id = $2", isPremium, userID)
	return err
}

func (r *UserRepository) ClaimRefs(userID int64) (rewardsGot int, err error) {
	var refsActivated, refsFrozen, refsPremiumActivatedNotClaimed, refsPremiumFrozenNotClaimed, refsActivatedNotClaimed, refsFrozenNotClaimed int

	row := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = true", userID)
	if err := row.Scan(&refsActivated); err != nil {
		return 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = false", userID)
	if err := row.Scan(&refsFrozen); err != nil {
		return 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsPremiumActivatedNotClaimed); err != nil {
		return 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = false AND ref_claimed = false", userID)
	if err := row.Scan(&refsPremiumFrozenNotClaimed); err != nil {
		return 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = false AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsActivatedNotClaimed); err != nil {
		return 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = false AND is_activated = false AND ref_claimed = false", userID)
	if err := row.Scan(&refsFrozenNotClaimed); err != nil {
		return 0, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	fmt.Printf("refsActivated: %d, refsFrozen: %d, refsPremiumActivatedNotClaimed: %d, refsPremiumFrozenNotClaimed: %d, refsActivatedNotClaimed: %d, refsFrozenNotClaimed: %d\n", refsActivated, refsFrozen, refsPremiumActivatedNotClaimed, refsPremiumFrozenNotClaimed, refsActivatedNotClaimed, refsFrozenNotClaimed)
	rewardsGot = refsActivatedNotClaimed*service.RefReward + refsPremiumActivatedNotClaimed*service.RefRewardPremium
	_, err = tx.Exec("UPDATE users SET ref_claimed = true WHERE referer = $1 AND is_activated = true", userID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec("UPDATE users SET point_balance = point_balance + $1 WHERE user_id = $2", rewardsGot, userID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return rewardsGot, nil
}
