package service

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/Corray333/therun_miniapp/internal/domains/city/types"
	user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"
)

// Errors
var (
	ErrMaxBuildingLevel = errors.New("max building level")
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetCity(ctx context.Context, userID int64) ([]types.BuildingPublic, error)
	GetBuilding(ctx context.Context, userID int64, buildingType types.BuildingType) (types.BuildingPublic, error)

	UpgradeBuilding(ctx context.Context, userID int64, buildingType types.BuildingType, buildingTime int64) error

	GetResources(ctx context.Context, userID int64) ([]types.Resource, error)

	ChangeBalances(ctx context.Context, userID int64, changes []user_types.BalanceChange) error
}

type CityService struct {
	repo repository
}

func New(repo repository) *CityService {
	return &CityService{
		repo: repo,
	}
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

		buildings[i].Img = getBuildingImage(buildings[i].Type, buildings[i].Level)

		buildingsMap[buildings[i].Type] = buildings[i]
	}

	return buildingsMap, nil
}

func getBuildingImage(buildingType types.BuildingType, level int) string {
	levelStr := strconv.Itoa(level)
	if level == 0 {
		levelStr = "1"
	}
	return os.Getenv("BASE_URL") + "/static/images/buildings/" + string(buildingType) + levelStr + ".png"
}

func (s *CityService) GetWarehouse(userID int64) (*types.WarehousePublic, error) {
	warehouseLevel, err := s.repo.GetBuilding(context.Background(), userID, types.BuildingWarehouse)
	if err != nil {
		return nil, err
	}
	warehouseLevel.Img = getBuildingImage(warehouseLevel.Type, warehouseLevel.Level)

	resources, err := s.repo.GetResources(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	for i := range resources {
		resources[i].Type = types.Resources[resources[i].Name].Type
	}

	var nextLevel *types.WarehouseLevel
	if warehouseLevel.Level < len(types.WarehouseLevels) {
		nextLevel = &types.WarehouseLevels[warehouseLevel.Level]
	}

	var currentLevel *types.WarehouseLevel
	if warehouseLevel.Level > 0 {
		currentLevel = &types.WarehouseLevels[warehouseLevel.Level-1]
	}

	var moreLevel *types.WarehouseLevel
	if warehouseLevel.Level == 0 && len(types.WarehouseLevels) > 1 {
		moreLevel = &types.WarehouseLevels[1]
	}

	warehouse := &types.WarehousePublic{
		BuildingPublic: warehouseLevel,
		Resources:      resources,
		CurrentLevel:   currentLevel,
		NextLevel:      nextLevel,
		MoreLevel:      moreLevel,
	}

	return warehouse, nil

}

func (s *CityService) UpgradeBuilding(userID int64, buildingType types.BuildingType) error {
	ctx, err := s.repo.Begin(context.Background())
	if err != nil {
		return err
	}
	defer s.repo.Rollback(ctx)

	building, err := s.repo.GetBuilding(ctx, userID, buildingType)
	if err != nil {
		return err
	}

	building.UpgradeCost = types.Buildings[building.Type].GetNextLevelCost(building.Level)
	if len(building.UpgradeCost) == 0 {
		return ErrMaxBuildingLevel
	}
	if err := s.repo.ChangeBalances(ctx, userID, building.UpgradeCost); err != nil {
		return err
	}

	if err := s.repo.UpgradeBuilding(ctx, userID, buildingType, types.Buildings[building.Type].GetNextLevelBuildTime(building.Level)); err != nil {
		return err
	}

	err = s.repo.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
