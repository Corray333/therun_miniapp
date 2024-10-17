package types

type Currency string

const (
	Point = Currency("point")
	Race  = Currency("race")

	RedKey   = Currency("red_key")
	BlueKey  = Currency("blue_key")
	GreenKey = Currency("green_key")
)

type User struct {
	ID       int64  `json:"id" db:"user_id"`
	Avatar   string `json:"avatar" db:"avatar"`
	Username string `json:"username,omitempty" db:"username"`
	InAppID  string `json:"-" db:"in_app_id"`

	PointBalance int `json:"pointBalance" db:"point_balance"`
	RaceBalance  int `json:"raceBalance" db:"race_balance"`

	RedKeyBalance   int `json:"red_keyBalance" db:"red_key_balance"`
	BlueKeyBalance  int `json:"blue_keyBalance" db:"blue_key_balance"`
	GreenKeyBalance int `json:"green_keyBalance" db:"green_key_balance"`

	LastClaim   int64 `json:"lastClaim" db:"last_claim"`
	FarmingFrom int64 `json:"farmingFrom" db:"farming_from"`
	MaxPoints   int   `json:"maxPoints" db:"max_points"`
	FarmingTime int   `json:"farmingTime" db:"farm_time"`

	RefCode    string `json:"refCode,omitempty" db:"ref_code"`
	Referer    *int64 `json:"referer,omitempty" db:"referer"`
	RefClaimed bool   `json:"refClaimed" db:"ref_claimed"`

	DailyCheckStreak int   `json:"dailyCheckStreak" db:"daily_check_streak"`
	DailyCheckLast   int64 `json:"dailyCheckLast" db:"daily_check_last"`

	IsPremium   bool `json:"isPremium" db:"is_premium"`
	IsActivated bool `json:"isActivated" db:"is_activated"`

	CurrentMiles  float64 `json:"currentMiles" db:"current_miles"`
	RaceStartTime int64   `json:"raceStartTime" db:"race_start_time"`
}

type Referal struct {
	Avatar    string `json:"avatar" db:"avatar"`
	Username  string `json:"username" db:"username"`
	IsPremium bool   `json:"isPremium" db:"is_premium"`
}

type BalanceChange struct {
	Currency Currency `json:"currency" db:"currency"`
	Amount   int      `json:"amount" db:"amount"`
}

var RefReward []BalanceChange = []BalanceChange{
	{
		Currency: Point,
		Amount:   1000,
	},
	{
		Currency: BlueKey,
		Amount:   1,
	},
}

var RefRewardPremium []BalanceChange = []BalanceChange{
	{
		Currency: Point,
		Amount:   3000,
	},
	{
		Currency: RedKey,
		Amount:   3,
	},
	{
		Currency: BlueKey,
		Amount:   2,
	},
}
