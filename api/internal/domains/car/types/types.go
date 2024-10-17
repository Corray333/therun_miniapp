package types

import (
	round_types "github.com/Corray333/therun_miniapp/internal/domains/round/types"
)

type Characteristic string

const (
	CharacteristicAcceleration Characteristic = "acceleration"
	CharacteristicSpeed        Characteristic = "speed"
	CharacteristicHendling     Characteristic = "hendling"
	CharacteristicBrakes       Characteristic = "brakes"
	CharacteristicStrength     Characteristic = "strength"
	CharacteristicFuel         Characteristic = "fuel"
)

type Car struct {
	UserID int64 `json:"-" db:"user_id"`
	ID     int64 `json:"id" db:"car_id"`
	IsMain bool  `json:"isMain" db:"is_main"`

	Element      round_types.Element `json:"element" db:"element"`
	Img          string              `json:"img" db:"-"`
	Acceleration int                 `json:"acceleration" db:"acceleration"`
	Hendling     int                 `json:"hendling" db:"hendling"`
	Brakes       int                 `json:"brakes" db:"brakes"`
	Strength     int                 `json:"strength" db:"strength"`
	Tank         int                 `json:"tank" db:"tank"`
	Fuel         float64             `json:"fuel" db:"fuel"`
	Health       float64             `json:"health" db:"health"`

	Modules []Module `json:"modules" db:"-"`
}

type Module struct {
	Characteristic Characteristic `json:"characteristic" db:"characteristic"`
	Boost          float64        `json:"boost" db:"boost"`
	Img            string         `json:"img" db:"-"`
	Name           string         `json:"name" db:"name"`
}

type RaceState struct {
	CurrentMiles float64 `json:"currentMiles" db:"miles"`
	StartTime    int64   `json:"startTime" db:"start_time"`
	Place        int     `json:"place" db:"place"`
}

var ElementEffects = map[round_types.Element]map[round_types.Element]int{
	round_types.ElementDesert: {
		round_types.ElementDesert: 100,
		round_types.ElementCity:   70,
		round_types.ElementTrack:  40,
	},
	round_types.ElementCity: {
		round_types.ElementDesert: 40,
		round_types.ElementCity:   100,
		round_types.ElementTrack:  70,
	},
	round_types.ElementTrack: {
		round_types.ElementDesert: 30,
		round_types.ElementCity:   80,
		round_types.ElementTrack:  100,
	},
}

var ElementEffectsFuel = map[round_types.Element]map[round_types.Element]int{
	round_types.ElementDesert: {
		round_types.ElementDesert: 100,
		round_types.ElementCity:   70,
		round_types.ElementTrack:  40,
	},
	round_types.ElementCity: {
		round_types.ElementDesert: 40,
		round_types.ElementCity:   100,
		round_types.ElementTrack:  70,
	},
	round_types.ElementTrack: {
		round_types.ElementDesert: 30,
		round_types.ElementCity:   80,
		round_types.ElementTrack:  100,
	},
}
