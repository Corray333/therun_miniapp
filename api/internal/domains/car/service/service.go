package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/car/types"
	round_types "github.com/Corray333/therun_miniapp/internal/domains/round/types"
)

type repository interface {
	BuyCar(ctx context.Context, car *types.Car, userID int64) error
	GetOwnedCars(ctx context.Context, userID int64) []types.Car
	GetMainCar(ctx context.Context, userID int64) (*types.Car, error)
	GetCarByID(ctx context.Context, carID int64) (*types.Car, error)

	PickCar(ttx context.Context, carID int64, userID int64) error
}

type userService interface {
}

type roundService interface {
	CurrentRound() *round_types.Round
}

type CarService struct {
	repo         repository
	userService  userService
	roundService roundService
}

func New(repo repository, userService userService, roundService roundService) *CarService {
	return &CarService{
		repo:         repo,
		userService:  userService,
		roundService: roundService,
	}
}

func (s *CarService) generateCar(element round_types.Element) *types.Car {
	car := types.Car{
		Element: element,
	}
	seed := int64(time.Now().UnixNano())
	r := rand.New(rand.NewSource(seed))

	min := 50
	max := 100

	car.Acceleration = r.Intn(max-min+1) + min
	car.Hendling = r.Intn(max-min+1) + min
	car.Brakes = r.Intn(max-min+1) + min
	car.Strength = r.Intn(max-min+1) + min

	return &car
}

func (s *CarService) countSpeed(car *types.Car) float64 {
	accelerationCoef := car.Acceleration / 100 * 25
	hendlingCoef := car.Hendling / 100 * 25
	brakesCoef := car.Brakes / 100 * 25
	strengthCoef := car.Strength / 100 * 25

	for _, module := range car.Modules {
		switch module.Characteristic {
		case types.CharacteristicAcceleration:
			accelerationCoef *= module.Boost
		case types.CharacteristicHendling:
			hendlingCoef *= module.Boost
		case types.CharacteristicBrakes:
			brakesCoef *= module.Boost
		case types.CharacteristicStrength:
			strengthCoef *= module.Boost
		}
	}

	round := s.roundService.CurrentRound()
	elementCoef := types.ElementEffects[car.Element][round.Element]

	return float64(accelerationCoef+hendlingCoef+brakesCoef+strengthCoef) * (float64(elementCoef) / 100)
}

func (s *CarService) countMiles(duration int, car *types.Car) float64 {
	speed := s.countSpeed(car)
	return speed * float64(duration/60/60)
}

func (s *CarService) countFuelWasting(duration int, car *types.Car) float64 {
	round := s.roundService.CurrentRound()
	speed := s.countSpeed(car)

	wasting := speed * 26 / 10 * float64(car.Fuel) / 100 / float64(types.ElementEffectsFuel[car.Element][round.Element]/100)

	miles := s.countMiles(duration, car)

	wasted := miles / wasting

	return wasted
}

func (s *CarService) GetAllCars(ctx context.Context) []types.Car {
	cars := make([]types.Car, 3)

	cars = append(cars, types.Car{
		Element: round_types.ElementDesert,
	}, types.Car{
		Element: round_types.ElementCity,
	}, types.Car{
		Element: round_types.ElementTrack,
	})

	return cars
}

func (s *CarService) BuyCar(ctx context.Context, element round_types.Element, userID int64) error {
	car := s.generateCar(element)

	if err := s.repo.BuyCar(ctx, car, userID); err != nil {
		return err
	}

	return nil
}

func (s *CarService) GetMainCar(ctx context.Context, userID int64) (*types.Car, error) {
	return s.repo.GetMainCar(ctx, userID)
}

func (s *CarService) GetCarByID(ctx context.Context, carID int64) (*types.Car, error) {
	return s.repo.GetCarByID(ctx, carID)
}

func (s *CarService) PickCar(ctx context.Context, carID int64, userID int64) error {
	return s.repo.PickCar(ctx, carID, userID)
}

func (s *CarService) GetOwnedCars(ctx context.Context, userID int64) []types.Car {
	return s.repo.GetOwnedCars(ctx, userID)
}

func (s *CarService) CurrentPoints(car *types.Car, state *types.RaceState) float64 {
	return s.countMiles(int(time.Now().Unix())-int(state.StartTime), car)
}

func (s *CarService) GetRace(ctx context.Context, userID int64) (race *types.RaceState, err error)
func (s *CarService) StartRace(ctx context.Context, userID int64) (race *types.RaceState, err error)
func (s *CarService) EndRace(ctx context.Context, userID int64) (race *types.RaceState, err error)
