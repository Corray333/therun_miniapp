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
			{Amount: 30, Probability: 20},
			{Amount: 20, Probability: 30},
			{Amount: 10, Probability: 50},
		},
		MinRewards: 10,
		MaxRewards: 30,
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
