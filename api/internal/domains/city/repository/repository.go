package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/city/types"
	user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

var (
	ErrInvalidTxType = errors.New("invalid transaction type")
)

type userRepository interface {
	ChangeBalances(ctx context.Context, userID int64, changes []user_types.BalanceChange) error
}

type CityRepository struct {
	db *sqlx.DB
	userRepository
}

func New(store *storage.Storage, userRepository userRepository) *CityRepository {
	return &CityRepository{
		db:             store.DB(),
		userRepository: userRepository,
	}
}

func (r *CityRepository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, storage.TxKey{}, tx), nil
}

func (r *CityRepository) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Commit()
}

func (r *CityRepository) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Rollback()
}

func (r *CityRepository) getTx(ctx context.Context) (tx *sqlx.Tx, isNew bool, err error) {
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

func (r *CityRepository) GetCity(ctx context.Context, userID int64) (buildings []types.BuildingPublic, err error) {
	var tx *sqlx.Tx
	txRaw := ctx.Value(storage.TxKey{})
	if txRaw != nil {
		var ok bool
		tx, ok = txRaw.(*sqlx.Tx)
		if !ok {
			slog.Error("invalid transaction type")
			return nil, nil
		}
	}
	if tx == nil {
		tx, err = r.db.BeginTxx(ctx, nil)
		if err != nil {
			slog.Error("failed to begin transaction: " + err.Error())
			return nil, err
		}
		defer tx.Rollback()
	}

	if err = tx.Select(&buildings, "SELECT type, level, state, state_until FROM buildings WHERE user_id = $1", userID); err != nil {
		slog.Error("failed to select buildings: " + err.Error())
		return nil, err
	}

	return buildings, nil

}

func (r *CityRepository) GetBuilding(ctx context.Context, userID int64, buildingType types.BuildingType) (building types.BuildingPublic, err error) {
	if err = r.db.Get(&building, "SELECT type, level, state, last_state_change, state_until FROM buildings WHERE type = $1 AND user_id = $2", buildingType, userID); err != nil {
		slog.Error("failed to get building: " + err.Error())
		return building, err
	}

	return building, nil
}

func (r *CityRepository) GetResources(ctx context.Context, userID int64) (resources []types.Resource, err error) {
	if err = r.db.Select(&resources, "SELECT name, amount FROM resources WHERE user_id = $1", userID); err != nil {
		slog.Error("failed to select resources: " + err.Error())
		return nil, err
	}

	return resources, nil
}

func (r *CityRepository) UpgradeBuilding(ctx context.Context, userID int64, buildingType types.BuildingType, buildingTime int64) error {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.Exec("UPDATE buildings SET level = level + 1, last_state_change = $3, state = $4, state_until = $5 WHERE type = $1 AND user_id = $2", buildingType, userID, time.Now().Unix(), types.BuildingStateBuild, time.Now().Unix()+buildingTime)
	if err != nil {
		slog.Error("failed to update building: " + err.Error())
		return err
	}

	if isNew {
		if err = tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}
	return nil
}
