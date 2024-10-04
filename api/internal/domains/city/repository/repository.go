package repository

import (
	"context"
	"log/slog"

	"github.com/Corray333/therun_miniapp/internal/domains/city/types"
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

	if err = tx.Select(&buildings, "SELECT type, level FROM buildings WHERE user_id = $1", userID); err != nil {
		slog.Error("failed to select buildings: " + err.Error())
		return nil, err
	}

	return buildings, nil

}
