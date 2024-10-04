package service

import (
	"context"
	"os"
	"strconv"

	"github.com/Corray333/therun_miniapp/internal/domains/city/types"
)

// Errors
var ()

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetCity(ctx context.Context, userID int64) ([]types.BuildingPublic, error)
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

func (s *CityService) GetCity(ctx context.Context, userID int64) (map[types.BuildingType]types.BuildingPublic, error) {
	buildings, err := s.repo.GetCity(ctx, userID)
	if err != nil {
		return nil, err
	}

	buildingsMap := make(map[types.BuildingType]types.BuildingPublic)

	for i := range buildings {
		building := types.Buildings[buildings[i].Type]
		if building != nil {
			buildings[i].UpgradeCost = building.GetNextLevelCost(buildings[i].Level)
		}
		level := strconv.Itoa(buildings[i].Level)
		if buildings[i].Level == 0 {
			level = "1"
		}
		buildings[i].Img = os.Getenv("BASE_URL") + "/static/images/buildings/" + string(buildings[i].Type) + level + ".png"

		buildingsMap[buildings[i].Type] = buildings[i]
	}

	return buildingsMap, nil
}

func (s *CityService) GetWarehouse(userID int64) (*types.WarehouseLevel, error) {
	return nil, nil
}

func (s *CityService) UpgradeBuilding(userID, buildingType types.BuildingType) error {
	return nil
}
