package repository

import (
	"fmt"
	"log/slog"

	"github.com/Corray333/therun_miniapp/internal/domains/cases/types"
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

func (t *CityRepository) HasCase(userID int64, caseType string) (hasCase bool, err error) {
	count := 0
	if err = t.db.Get(&count, "SELECT COUNT(*) FROM user_case WHERE user_id = $1 AND case_type = $2", userID, caseType); err != nil {
		slog.Error("error while checking case" + err.Error())
		return false, err
	}
	return count > 0, nil
}

func (t *CityRepository) OpenCase(userID int64, case_type *types.Case, reward *types.Reward) error {
	tx, err := t.db.Beginx()
	if err != nil {
		slog.Error("error while opening user_case transaction" + err.Error())
		return err
	}

	defer tx.Rollback()

	cost := ""
	for _, keys := range case_type.Keys {
		cost += fmt.Sprintf("%s_balance = %s_balance - %d,", keys.Type, keys.Type, keys.Amount)
	}
	if len(cost) == 0 {
		return fmt.Errorf("no keys in case")
	}
	cost = cost[:len(cost)-1]
	fmt.Println(fmt.Sprintf("UPDATE users SET %s_balance = %s_balance + $1, %s WHERE user_id = $2", case_type.RewardType, case_type.RewardType, cost))

	if _, err = tx.Exec(fmt.Sprintf("UPDATE users SET %s_balance = %s_balance + $1, %s WHERE user_id = $2", case_type.RewardType, case_type.RewardType, cost), reward.Amount, userID); err != nil {
		slog.Error("error while updating user balance" + err.Error())
		return err
	}

	// TODO: remove case or do something else
	return tx.Commit()
}

func (t *CityRepository) GetCases(userID int64) (cases []types.Case, err error) {
	if err = t.db.Select(&cases, "SELECT case_type FROM user_case WHERE user_id = $1", userID); err != nil {
		slog.Error("error while getting cases" + err.Error())
		return nil, err
	}
	return cases, nil
}
