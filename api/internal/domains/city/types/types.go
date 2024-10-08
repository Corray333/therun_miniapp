package types

import (
	"time"

	user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"
)

type Resource struct {
	Name   string       `json:"name"`
	Type   ResourceType `json:"type"`
	Amount int          `json:"amount"`
}

var Resources = map[string]Resource{
	"titan": {
		Name:   "titan",
		Type:   ResourceTypeOre,
		Amount: 0,
	},
	"quartz": {
		Name:   "quartz",
		Type:   ResourceTypeMineral,
		Amount: 0,
	},
}

type ResourceType string
type BuildingType string
type BuildingState string

const (
	ResourceTypeMetal   = ResourceType("metal")
	ResourceTypeWood    = ResourceType("wood")
	ResourceTypeMineral = ResourceType("mineral")
	ResourceTypeOre     = ResourceType("ore")
	ResourceTypeCrystal = ResourceType("crystal")
)

const (
	BuildingStateIdle  = BuildingState("idle")
	BuildingStateBuild = BuildingState("build")
	BuildingStateWork  = BuildingState("work")
)

const (
	BuildingWarehouse = BuildingType("warehouse")
	BuildingMine      = BuildingType("mine")
	BuildingFabric    = BuildingType("fabric")
	BuildingForest    = BuildingType("forest")
)

type Requirement struct {
	Type  BuildingType `json:"type"`
	Level int          `json:"level"`
}

type WarehouseLevel struct {
	Capacity         int                        `json:"capacity"`
	Cost             []user_types.BalanceChange `json:"cost"`
	Requirements     []Requirement              `json:"requirements"`
	BuildingDuration time.Duration              `json:"buildingDuration"`
}

type Warehouse []WarehouseLevel

func (w Warehouse) GetNextLevelCost(level int) []user_types.BalanceChange {
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
		BuildingDuration: 10 * 60, // 10 minutes
		Cost: []user_types.BalanceChange{
			{
				Currency: user_types.Point,
				Amount:   -1000,
			},
			{
				Currency: user_types.BlueKey,
				Amount:   -1,
			},
		},
		Requirements: nil,
	},
	WarehouseLevel{
		Capacity:         2000,
		BuildingDuration: 2 * 60 * 60, // 2 hours
		Cost: []user_types.BalanceChange{
			{
				Currency: user_types.Point,
				Amount:   -2000,
			},
			{
				Currency: user_types.BlueKey,
				Amount:   -2,
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
	GetNextLevelCost(level int) []user_types.BalanceChange
}

type BuildingPublic struct {
	Img             string                     `json:"img"`
	Type            BuildingType               `json:"type" db:"type"`
	Level           int                        `json:"level" db:"level"`
	State           BuildingState              `json:"state" db:"state"`
	LastStateChange int64                      `json:"lastStateChange" db:"last_state_change"`
	UpgradeCost     []user_types.BalanceChange `json:"upgradeCost,omitempty"`
}

type WarehousePublic struct {
	BuildingPublic

	Resources    []Resource      `json:"resources"`
	CurrentLevel *WarehouseLevel `json:"currentLevel"`
	NextLevel    *WarehouseLevel `json:"nextLevel"`
}

var Buildings = map[BuildingType]Building{
	BuildingWarehouse: WarehouseLevels,
}
