package repository

import (
	"log/slog"

	"github.com/Corray333/therun_miniapp/internal/domains/round/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type RoundRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *RoundRepository {
	return &RoundRepository{
		db: store.DB(),
	}
}

func (r *RoundRepository) SetRound(round *types.Round) error {
	_, err := r.db.Exec("INSERT INTO rounds (round_id, end_time, element) VALUES ($1, $2, $3) ON CONFLICT (round_id) DO UPDATE SET end_time = $2", round.ID, round.EndTime, round.Element)
	if err != nil {
		slog.Error("error setting round: " + err.Error())
		return err
	}
	return nil
}

func (r *RoundRepository) GetRoundElement(roundID int) (types.Element, error) {
	var element types.Element
	err := r.db.Get(&element, "SELECT element FROM rounds WHERE round_id = $1", roundID)
	if err != nil {
		slog.Error("error getting round element: " + err.Error())
		return types.Element(""), err
	}
	return element, nil
}
