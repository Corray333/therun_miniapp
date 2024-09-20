package types

import "encoding/json"

// Task types
const (
	TaskTypeExternal    = "external"
	TaskTypeTg          = "tg"
	TaskTypeUncheckable = "uncheckable"
)

type Tanslation struct {
	Lang string `json:"lang" db:"lang"`
	Text string `json:"text" db:"text"`
}

type Task struct {
	ID           int             `json:"id" db:"task_id"`
	Description  string          `json:"description" db:"description"`
	Type         string          `json:"type" db:"type"`
	Link         string          `json:"link" db:"link"`
	ExpireAt     int64           `json:"expireAt" db:"expire_at"`
	PointsReward int             `json:"pointsReward" db:"points_reward"`
	KeysReward   int             `json:"keysReward" db:"keys_reward"`
	RaceReward   int             `json:"raceReward" db:"race_reward"`
	Data         json.RawMessage `json:"data" db:"data"`
	Icon         string          `json:"icon" db:"icon"`
}
