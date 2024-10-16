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
	Element      round_types.Element `json:"element" db:"element"`
	Img          string              `json:"img" db:"-"`
	Acceleration int                 `json:"acceleration" db:"acceleration"`
	Hendling     int                 `json:"hendling" db:"hendling"`
	Brakes       int                 `json:"brakes" db:"brakes"`
	Strength     int                 `json:"strength" db:"strength"`
	Fuel         int                 `json:"fuel" db:"fuel"`

	Modules []Module `json:"modules" db:"-"`
}

type Module struct {
	Characteristic Characteristic `json:"characteristic" db:"characteristic"`
	Boost          int            `json:"boost" db:"boost"`
	Img            string         `json:"img" db:"-"`
	Name           string         `json:"name" db:"name"`
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
