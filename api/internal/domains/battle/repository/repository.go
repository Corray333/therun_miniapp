package repository

import (
	"log/slog"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/service"
	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/jmoiron/sqlx"
)

type BattleRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *BattleRepository {
	return &BattleRepository{
		db: store.DB(),
	}
}

type battleInternal struct {
	ID               int     `json:"id" db:"battle_id"`
	RoundID          int     `json:"roundID" db:"round_id"`
	UserID           int     `json:"userID" db:"user_id"`
	UserUsername     string  `json:"userUsername" db:"user_username"`
	UserPhoto        string  `json:"userPhoto" db:"user_photo"`
	UserMiles        float64 `json:"userResult" db:"user_miles"`
	OpponentID       int     `json:"opponentID" db:"opponent_id"`
	OpponentUsername string  `json:"opponentUsername" db:"opponent_username"`
	OpponentPhoto    string  `json:"opponentPhoto" db:"opponent_photo"`
	OpponentMiles    float64 `json:"opponentResult" db:"opponent_miles"`
	Pick             *int    `json:"pick" db:"pick"`
}

func (r *BattleRepository) GetBattles(round int, userID int64) ([]types.Battle, error) {
	battles := []battleInternal{}
	err := r.db.Select(&battles, "SELECT battles.*, bets.pick FROM battles LEFT JOIN bets ON battles.battle_id = bets.battle_id AND bets.user_id = $2 WHERE round_id = $1", round, userID)
	if err != nil {
		return nil, err
	}

	var result []types.Battle
	for _, battle := range battles {
		pick := 0
		if battle.Pick != nil {
			pick = *battle.Pick
		}
		result = append(result, types.Battle{
			ID:      battle.ID,
			RoundID: battle.RoundID,
			User: types.User{
				ID:       battle.UserID,
				Username: battle.UserUsername,
				Photo:    battle.UserPhoto,
			},
			Opponent: types.Opponent{
				ID:       battle.OpponentID,
				Username: battle.OpponentUsername,
				Photo:    battle.OpponentPhoto,
			},
			UserMiles:     battle.UserMiles,
			OpponentMiles: battle.OpponentMiles,
			Pick:          pick,
		})
	}

	return result, nil
}

func (r *BattleRepository) SetBattles(battles []types.Battle) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, battle := range battles {
		_, err := tx.Exec(`
			INSERT INTO battles (battle_id, round_id, user_id, user_username, user_photo, user_miles, opponent_id, opponent_username, opponent_photo, opponent_miles) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
			ON CONFLICT (battle_id) 
			DO UPDATE SET 
				round_id = EXCLUDED.round_id,
				user_id = EXCLUDED.user_id,
				user_username = EXCLUDED.user_username,
				user_photo = EXCLUDED.user_photo,
				user_miles = EXCLUDED.user_miles,
				opponent_id = EXCLUDED.opponent_id,
				opponent_username = EXCLUDED.opponent_username,
				opponent_photo = EXCLUDED.opponent_photo,
				opponent_miles = EXCLUDED.opponent_miles
		`, battle.ID, battle.RoundID, battle.User.ID, battle.User.Username, battle.User.Photo, battle.UserMiles, battle.Opponent.ID, battle.Opponent.Username, battle.Opponent.Photo, battle.OpponentMiles)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *BattleRepository) GetBattle(battleID int) (*types.Battle, error) {
	battle := battleInternal{}
	err := r.db.Get(&battle, "SELECT * FROM battles WHERE battle_id = $1", battleID)
	if err != nil {
		return nil, err
	}

	pick := 0
	if battle.Pick != nil {
		pick = *battle.Pick
	}
	return &types.Battle{
		ID:      battle.ID,
		RoundID: battle.RoundID,
		User: types.User{
			ID:       battle.UserID,
			Username: battle.UserUsername,
			Photo:    battle.UserPhoto,
		},
		Opponent: types.Opponent{
			ID:       battle.OpponentID,
			Username: battle.OpponentUsername,
			Photo:    battle.OpponentPhoto,
		},
		UserMiles:     battle.UserMiles,
		OpponentMiles: battle.OpponentMiles,
		Pick:          pick,
	}, nil
}

func (r *BattleRepository) MakeBet(userID int64, battleID, pick int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO bets (user_id, battle_id, pick) VALUES ($1, $2, $3) ON CONFLICT (user_id, battle_id) DO UPDATE SET pick = $3", userID, battleID, pick)
	if err != nil {
		slog.Error("error making bet: " + err.Error())
		return err
	}

	if _, err := tx.Exec("UPDATE users SET point_balance = point_balance - $2 WHERE user_id = $1", userID, service.BetAmount); err != nil {
		slog.Error("error updating user balance: " + err.Error())
		return err
	}

	return tx.Commit()

}

func (r *BattleRepository) ProcessBets(roundID int, keyReward int) error {
	_, err := r.db.Exec(`WITH 
		round_battles AS (
			SELECT battle_id, user_miles, opponent_miles FROM battles WHERE round_id = $1
		), 
		winners AS (
			SELECT user_id from round_battles NATURAL JOIN bets WHERE user_miles > opponent_miles AND pick = 1 OR opponent_miles > user_miles AND pick = 2 
		)
		UPDATE users SET key_balance = key_balance + $2 WHERE user_id IN (SELECT user_id FROM winners)`, roundID, keyReward)
	if err != nil {
		slog.Error("error processing bets: " + err.Error())
		return err
	}

	return nil
}
