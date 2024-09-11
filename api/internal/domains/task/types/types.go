package types

// Task types
const (
	TaskTypeExternal    = "external"
	TaskTypeTg          = "tg"
	TaskTypeUncheckable = "uncheckable"
)

type Task struct {
	ID           int    `json:"id" db:"task_id"`
	Description  string `json:"description" db:"description"`
	Type         string `json:"type" db:"type"`
	Link         string `json:"link" db:"link"`
	ExpireAt     string `json:"expire_at" db:"expire_at"`
	PointsReward int    `json:"points_reward" db:"points_reward"`
	KeysReward   int    `json:"keys_reward" db:"keys_reward"`
	RaceReward   int    `json:"race_reward" db:"race_reward"`
}
