package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/car/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

var (
	ErrAlreadyHasCar      = errors.New("user already has a car")
	ErrRaceAlreadyStarted = errors.New("race already started")
	ErrNoCar              = errors.New("user has no car")
	ErrInvalidTxType      = errors.New("invalid transaction type")
)

type CarRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *CarRepository {
	return &CarRepository{
		db: store.DB(),
	}
}

func (r *CarRepository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, storage.TxKey{}, tx), nil
}

func (r *CarRepository) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Commit()
}

func (r *CarRepository) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Rollback()
}

func (r *CarRepository) getTx(ctx context.Context) (tx *sqlx.Tx, isNew bool, err error) {
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

func (r *CarRepository) GetMainCar(ctx context.Context, userID int64) (*types.Car, error) {
	var car types.Car
	err := r.db.Get(&car, "SELECT * FROM cars WHERE user_id = $1 AND is_main = true", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("error while getting main car: " + err.Error())
		return nil, err
	}

	return &car, nil
}

func (r *CarRepository) BuyCar(ctx context.Context, car *types.Car, userID int64) error {

	oldCar, err := r.GetMainCar(ctx, userID)
	if err != nil {
		return err
	}

	if oldCar != nil {
		return ErrAlreadyHasCar
	}

	_, err = r.db.Exec("INSERT INTO cars (user_id, element, acceleration, hendling, brakes, strength, tank, fuel, health, is_main) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", userID, car.Element, car.Acceleration, car.Hendling, car.Brakes, car.Strength, car.Tank, car.Fuel, car.Health, car.IsMain)
	if err != nil {
		slog.Error("error while choosing car: " + err.Error())
		return err
	}

	return nil
}

func (r *CarRepository) GetCarByID(ctx context.Context, carID int64) (*types.Car, error) {
	var car types.Car
	err := r.db.Get(&car, "SELECT * FROM cars WHERE car_id = $1", carID)
	if err != nil {
		slog.Error("error while getting car by id: " + err.Error())
		return nil, err
	}

	return &car, nil
}

func (r *CarRepository) PickCar(ctx context.Context, carID int64, userID int64) error {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err = tx.Exec("UPDATE cars SET is_main = false WHERE user_id = $1", userID); err != nil {
		slog.Error("error while picking car: " + err.Error())
		return err
	}

	_, err = tx.Exec("UPDATE cars SET is_main = true WHERE car_id = $1 AND user_id = $2", carID, userID)
	if err != nil {
		slog.Error("error while picking car: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}

func (r *CarRepository) GetOwnedCars(ctx context.Context, userID int64) []types.Car {
	var cars []types.Car
	err := r.db.Select(&cars, "SELECT * FROM cars WHERE user_id = $1", userID)
	if err != nil {
		slog.Error("error while getting owned cars: " + err.Error())
		return nil
	}

	return cars
}

func (r *CarRepository) GetRaceState(ctx context.Context, userID int64, roundID int) (*types.RaceState, error) {
	var raceState types.RaceState

	if err := r.db.Get(&raceState, "SELECT start_time, miles FROM races WHERE user_id = $1 AND round_id = $2", userID, roundID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("error while getting race state: " + err.Error())
		return nil, err
	}

	return &raceState, nil
}

func (r *CarRepository) StartRace(ctx context.Context, userID int64, roundID int) (*types.RaceState, error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}
	if isNew {
		defer tx.Rollback()
	}

	car, err := r.GetMainCar(ctx, userID)
	if err != nil {
		return nil, err
	}

	if car == nil {
		return nil, ErrNoCar
	}

	state, err := r.GetRaceState(ctx, userID, roundID)
	if err != nil {
		return nil, err
	}

	if state == nil {
		state = &types.RaceState{
			CurrentMiles: 0,
			StartTime:    0,
		}
	}

	if state.StartTime != 0 {
		return nil, ErrRaceAlreadyStarted
	}

	startTime := time.Now().Unix()

	_, err = tx.Exec("INSERT INTO races (user_id, round_id, start_time) VALUES ($1, $2, $3) ON CONFLICT (user_id, round_id) DO UPDATE SET start_time = $3", userID, roundID, startTime)
	if err != nil {
		slog.Error("error starting race: " + err.Error())
		return nil, err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return nil, err
		}
	}

	state.StartTime = startTime

	return state, nil
}

func (r *CarRepository) EndRace(ctx context.Context, userID int64, roundID int, miles float64, fuelWasted, healthWasted float64) error {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("UPDATE races SET miles = $1, start_time = 0 WHERE user_id = $2 AND round_id = $3", miles, userID, roundID); err != nil {
		slog.Error("error ending race: " + err.Error())
		return err
	}

	if _, err := tx.Exec("UPDATE cars SET fuel = fuel - $1, health = health - $2 WHERE user_id = $3 AND is_main = true", fuelWasted, healthWasted, userID); err != nil {
		slog.Error("error ending race: error updating car: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}
