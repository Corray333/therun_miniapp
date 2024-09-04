package types

type ContextKey string

type User struct {
	ID           int64  `json:"id,omitempty" db:"user_id"`
	Avatar       string `json:"avatar" db:"avatar"`
	Username     string `json:"username,omitempty" db:"username"`
	InAppID      string `json:"-" db:"in_app_id"`
	PointBalance int    `json:"pointBalance,omitempty" db:"point_balance"`
	RaceBalance  int    `json:"raceBalance,omitempty" db:"race_balance"`
	KeyBalance   int    `json:"keyBalance,omitempty" db:"key_balance"`
	LastClaim    int64  `json:"lastClaim,omitempty" db:"last_claim"`
	MaxPoints    int    `json:"maxPoints,omitempty" db:"max_points"`
	FarmTime     int    `json:"farmTime,omitempty" db:"farm_time"`
	RefCode      string `json:"refCode,omitempty" db:"ref_code"`
	Referer      *int64 `json:"referer,omitempty" db:"referer"`
}

type Referal struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}
