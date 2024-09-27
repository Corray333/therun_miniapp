package types

type Battle struct {
	ID       int  `json:"id" db:"battle_id"`
	RoundID  int  `json:"roundID" db:"round_id"`
	User     User `json:"user"`
	Opponent User `json:"opponent"`
}

type User struct {
	ID       int     `json:"id" db:"user_id"`
	Username string  `json:"username" db:"username"`
	Photo    string  `json:"photo" db:"photo"`
	Miles    float64 `json:"miles" db:"miles"`
}

type Round struct {
	ID      int      `json:"id"`
	EndTime int      `json:"endTime"`
	Battles []Battle `json:"battles"`
}

type Bet struct {
	BattleID int `json:"battleID" db:"battle_id"`
	UserID   int `json:"userID" db:"user_id"`
	Pick     int `json:"pick" db:"pick"`
}
