package types

import (
	"time"

	user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"
)

type Resource struct {
	Name   string
	Type   string
	Amount int
}

type ResourceType string
type BuildingType string

const (
	ResourceTitan  = ResourceType("titan")
	ResourceQuartz = ResourceType("quartz")
)

const (
	BuildingWarehouse = BuildingType("warehouse")
	BuildingMine      = BuildingType("mine")
	BuildingFabric    = BuildingType("fabric")
	BuildingForest    = BuildingType("forest")
)

type Cost struct {
	Currency user_types.Currency `json:"currency"`
	Amount   int                 `json:"amount"`
}

type Requirement struct {
	Type  BuildingType
	Level int
}

type WarehouseLevel struct {
	Capacity         int
	Resources        []ResourceType
	Cost             []Cost
	Requirements     []Requirement
	BuildingDuration time.Duration
}

type Warehouse []WarehouseLevel

func (w Warehouse) GetNextLevelCost(level int) []Cost {
	if level < 0 {
		return nil
	}
	if level >= len(w) {
		return nil
	}
	return w[level].Cost
}

var WarehouseLevels = Warehouse{
	WarehouseLevel{
		Capacity:         1000,
		BuildingDuration: 10 * time.Minute,
		Resources: []ResourceType{
			ResourceTitan,
		},
		Cost: []Cost{
			{
				Currency: user_types.Point,
				Amount:   1000,
			},
			{
				Currency: user_types.BlueKey,
				Amount:   1,
			},
		},
		Requirements: nil,
	},
	WarehouseLevel{
		Capacity:         2000,
		BuildingDuration: 2 * time.Hour,
		Resources: []ResourceType{
			ResourceTitan,
			ResourceQuartz,
		},
		Cost: []Cost{
			{
				Currency: user_types.Point,
				Amount:   2000,
			},
			{
				Currency: user_types.BlueKey,
				Amount:   2,
			},
		},
		Requirements: []Requirement{
			{
				Type:  BuildingWarehouse,
				Level: 1,
			},
		},
	},
}

type Building interface {
	GetNextLevelCost(level int) []Cost
}

type BuildingPublic struct {
	Img         string       `json:"img"`
	Type        BuildingType `json:"type"`
	Level       int          `json:"level"`
	UpgradeCost []Cost       `json:"upgradeCost"`
}

var Buildings = map[BuildingType]Building{
	BuildingWarehouse: WarehouseLevels,
}
