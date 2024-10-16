package types

type Battle struct {
	ID            int      `json:"id" db:"battle_id"`
	RoundID       int      `json:"roundID" db:"round_id"`
	User          User     `json:"user"`
	Opponent      Opponent `json:"opponent"`
	Pick          int      `json:"pick"`
	UserMiles     float64  `json:"userResult" db:"user_miles"`
	OpponentMiles float64  `json:"opponentResult" db:"opponent_miles"`
	Status        string   `json:"status"`
}

type User struct {
	ID       int    `json:"id" db:"user_id"`
	Username string `json:"username" db:"user_username"`
	Photo    string `json:"photo" db:"user_photo"`
}

type Opponent struct {
	ID       int    `json:"id" db:"opponent_id"`
	Username string `json:"username" db:"opponent_username"`
	Photo    string `json:"photo" db:"opponent_photo"`
}

type Bet struct {
	BattleID int `json:"battleID" db:"battle_id"`
	UserID   int `json:"userID" db:"user_id"`
	Pick     int `json:"pick" db:"pick"`
}
