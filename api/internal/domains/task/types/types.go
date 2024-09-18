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
	ExpireAt     int64           `json:"expire_at" db:"expire_at"`
	PointsReward int             `json:"points_reward" db:"points_reward"`
	KeysReward   int             `json:"keys_reward" db:"keys_reward"`
	RaceReward   int             `json:"race_reward" db:"race_reward"`
	Data         json.RawMessage `json:"data" db:"data"`
}

var example1 = Task{
	ID:           1,
	Description:  "Subscribe to HAPI telegram channel",
	Type:         TaskTypeTg,
	Link:         "https://t.me/hapi_ann",
	ExpireAt:     1610000000,
	PointsReward: 500,
	KeysReward:   0,
	RaceReward:   0,
	Data:         json.RawMessage(`{"tg_channel_id": 1610000000}`),
}
