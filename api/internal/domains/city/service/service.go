package service

import "github.com/Corray333/therun_miniapp/internal/domains/city/types"

// Errors
var ()

type repository interface {
	GetCity(userID int64) ([]types.BuildingPublic, error)
}

type CityService struct {
	repo repository
}

func New(repo repository) *CityService {
	return &CityService{
		repo: repo,
	}
}

func (s *CityService) getUserResources(userID int64) ([]types.Resource, error) {
	return nil, nil
}

func (s *CityService) GetCity(userID int64) ([]types.BuildingPublic, error) {
	buildings, err := s.repo.GetCity(userID)
	if err != nil {
		return nil, err
	}

	for i := range buildings {
		buildings[i].UpgradeCost = types.Buildings[types.BuildingType(buildings[i].Name)].GetLevelCost(buildings[i].Level + 1)
	}

	return buildings, nil
}

func (s *CityService) GetWarehouse(userID int64) (*types.WarehouseLevel, error) {
	return nil, nil
}

func (s *CityService) UpgradeBuilding(userID, buildingType types.BuildingType) error {
	return nil
}
