package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

var (
	ErrInvalidTxType = errors.New("invalid transaction type")
)

type UserRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *UserRepository {
	return &UserRepository{
		db: store.DB(),
	}
}

func (r *UserRepository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, storage.TxKey{}, tx), nil
}

func (r *UserRepository) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Commit()
}

func (r *UserRepository) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Rollback()
}

func (r *UserRepository) getTx(ctx context.Context) (tx *sqlx.Tx, isNew bool, err error) {
	txRaw := ctx.Value(storage.TxKey{})
	if txRaw != nil {
		var ok bool
		tx, ok = txRaw.(*sqlx.Tx)
		if !ok {
			slog.Error("invalid transaction type")
			return nil, false, ErrInvalidTxType
		}
	}
	if tx == nil {
		tx, err = r.db.BeginTxx(ctx, nil)
		if err != nil {
			slog.Error("failed to begin transaction: " + err.Error())
			return nil, false, err
		}

		return tx, true, nil
	}

	return tx, false, nil
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
        red_key_balance = :red_key_balance, 
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
	// Trigger inserts buildings, resources etc. for new users
	_, err := r.db.NamedExec("INSERT INTO users (user_id, username, avatar, in_app_id, point_balance, race_balance, red_key_balance, last_claim, max_points, farm_time, ref_code, referer, is_premium, is_activated) VALUES (:user_id, :username, :avatar, 0, :point_balance, :race_balance, :red_key_balance, :last_claim, :max_points, :farm_time, :ref_code, :referer, :is_premium, :is_activated)", user)
	return err
}

// TODO: check if this method works
func (r *UserRepository) CheckIfRefCodeExists(refCode string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE ref_code = $1)", refCode)
	return exists, err
}

// TODO: add pagination
func (r *UserRepository) ListActivatedReferals(userID int64) ([]types.Referal, error) {
	referals := []types.Referal{}
	err := r.db.Select(&referals, "SELECT avatar, username, is_premium FROM users WHERE referer = $1 AND is_activated = true", userID)
	return referals, err
}

// TODO: add pagination
func (r *UserRepository) ListNotActivatedReferals(userID int64) ([]types.Referal, error) {
	referals := []types.Referal{}
	err := r.db.Select(&referals, "SELECT avatar, username, is_premium FROM users WHERE referer = $1 AND is_activated = false", userID)
	return referals, err
}

func (r *UserRepository) CountReferals(userID int64) (refsActivated, refsFrozenTotal, refsFrozen, refsPremiumFrozen, refsPremiumActivatedNotClaimed, refsActivatedNotClaimed int, err error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = true", userID)
	if err := row.Scan(&refsActivated); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = false", userID)
	if err := row.Scan(&refsFrozenTotal); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_activated = false AND is_premium = false", userID)
	if err := row.Scan(&refsFrozen); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsPremiumActivatedNotClaimed); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = false AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsActivatedNotClaimed); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = false", userID)
	if err := row.Scan(&refsPremiumFrozen); err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	return refsActivated, refsFrozenTotal, refsFrozen, refsPremiumFrozen, refsPremiumActivatedNotClaimed, refsActivatedNotClaimed, nil
}

func (r *UserRepository) ChangeBalances(ctx context.Context, userID int64, changes []types.BalanceChange) error {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	for _, change := range changes {
		fmt.Println("UPDATE users SET " + string(change.Currency) + "_balance = " + string(change.Currency) + "_balance + $1 WHERE user_id = $2")
		_, err := tx.Exec("UPDATE users SET "+string(change.Currency)+"_balance = "+string(change.Currency)+"_balance + $1 WHERE user_id = $2", change.Amount, userID)
		if err != nil {
			return err
		}
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}

func (r *UserRepository) SetPremium(userID int64, isPremium bool) error {
	_, err := r.db.Exec("UPDATE users SET is_premium = $1 WHERE user_id = $2", isPremium, userID)
	return err
}

func (r *UserRepository) ClaimRefs(ctx context.Context, userID int64) (refsPremiumActivatedNotClaimed, refsActivatedNotClaimed int, err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return 0, 0, err
	}
	if isNew {
		defer tx.Rollback()
	}

	row := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = true AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsPremiumActivatedNotClaimed); err != nil {
		return 0, 0, err
	}

	row = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE referer = $1 AND is_premium = false AND is_activated = true AND ref_claimed = false", userID)
	if err := row.Scan(&refsActivatedNotClaimed); err != nil {
		return 0, 0, err
	}

	_, err = tx.Exec("UPDATE users SET ref_claimed = true WHERE referer = $1 AND is_activated = true", userID)
	if err != nil {
		return 0, 0, err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return 0, 0, err
		}
	}

	return
}

func (r *UserRepository) ActivateUser(userID int64) error {

	_, err := r.db.Exec("UPDATE users SET is_activated = true WHERE user_id = $1", userID)
	if err != nil {
		slog.Error("error activating user: " + err.Error())
		return err
	}

	return nil
}
