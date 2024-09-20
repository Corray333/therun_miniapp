package types

type User struct {
	ID       int64  `json:"id" db:"user_id"`
	Avatar   string `json:"avatar" db:"avatar"`
	Username string `json:"username,omitempty" db:"username"`
	InAppID  string `json:"-" db:"in_app_id"`

	PointBalance int `json:"pointBalance" db:"point_balance"`
	RaceBalance  int `json:"raceBalance" db:"race_balance"`
	KeyBalance   int `json:"keyBalance" db:"key_balance"`

	LastClaim   int64 `json:"lastClaim" db:"last_claim"`
	FarmingFrom int64 `json:"farmingFrom" db:"farming_from"`
	MaxPoints   int   `json:"maxPoints" db:"max_points"`
	FarmingTime int   `json:"farmingTime" db:"farm_time"`

	RefCode     string `json:"refCode,omitempty" db:"ref_code"`
	Referer     *int64 `json:"referer,omitempty" db:"referer"`
	RefsClaimed int    `json:"refsClaimed" db:"refs_claimed"`

	DailyCheckStreak int   `json:"dailyCheckStreak" db:"daily_check_streak"`
	DailyCheckLast   int64 `json:"dailyCheckLast" db:"daily_check_last"`
}

type Referal struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}
