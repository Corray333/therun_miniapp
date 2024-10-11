package types

const (
	CaseTypeRed = "red"
)

const (
	KeysTypeRed = "red_key"
)

const (
	RewardTypeRace = "race"
)

var Cases map[string]Case = map[string]Case{
	CaseTypeRed: {
		Type: CaseTypeRed,
		Keys: []Keys{
			{
				Type:   KeysTypeRed,
				Amount: 3,
			},
		},
		RewardType: RewardTypeRace,
		Rewards: []Reward{
			{Amount: 800, Probability: 10},
			{Amount: 700, Probability: 10},
			{Amount: 600, Probability: 20},
			{Amount: 500, Probability: 30},
			{Amount: 400, Probability: 30},
		},
		MinRewards: 400,
		MaxRewards: 800,
	},
}

type Reward struct {
	Type        string `json:"type"`
	Amount      int    `json:"amount"`
	Probability int    `json:"probability"`
}

type Case struct {
	Type       string   `json:"type" db:"case_type"`
	Keys       []Keys   `json:"keys"`
	RewardType string   `json:"rewardType"`
	Rewards    []Reward `json:"-"`

	MinRewards int `json:"min_rewards"`
	MaxRewards int `json:"max_rewards"`
}

type Keys struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}
