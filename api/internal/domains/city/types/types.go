package types

import user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"

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
	Type   user_types.Currency
	Amount int
}

type Requirement struct {
	Type  BuildingType
	Level int
}

type WarehouseLevel struct {
	Capacity     int
	Resources    []ResourceType
	Cost         []Cost
	Requirements []Requirement
}

type Warehouse []WarehouseLevel

func (w Warehouse) GetLevelCost(level int) []Cost {
	if level == 0 {
		return nil
	}
	if level >= len(w) {
		return nil
	}
	return w[level].Cost
}

var WarehouseLevels = Warehouse{
	WarehouseLevel{
		Capacity: 1000,
		Resources: []ResourceType{
			ResourceTitan,
		},
		Cost: []Cost{
			{
				Type:   user_types.Point,
				Amount: 1000,
			},
			{
				Type:   user_types.BlueKey,
				Amount: 1,
			},
		},
		Requirements: nil,
	},
	WarehouseLevel{
		Capacity: 2000,
		Resources: []ResourceType{
			ResourceTitan,
			ResourceQuartz,
		},
		Cost: []Cost{
			{
				Type:   user_types.Point,
				Amount: 2000,
			},
			{
				Type:   user_types.BlueKey,
				Amount: 2,
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
	GetLevelCost(level int) []Cost
}

type BuildingPublic struct {
	Name        string
	Level       int
	UpgradeCost []Cost
}

var Buildings = map[BuildingType]Building{
	BuildingWarehouse: WarehouseLevels,
}
