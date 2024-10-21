package service

import (
	"context"
	"fmt"
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

	GetRaceState(ctx context.Context, userId int64, roundID int) (*types.RaceState, error)
	StartRace(ctx context.Context, userID int64, roundID int) (*types.RaceState, error)
	EndRace(ctx context.Context, userID int64, roundID int, miles float64, fuelWasted, healthWasted float64) error

	GetModulesOfUser(ctx context.Context, carID int64) ([]types.Module, error)

	BuyFuel(ctx context.Context, userID int64, cost int) error
	BuyHealth(ctx context.Context, userID int64, cost int) error
}

type userService interface {
}

type roundService interface {
	GetRound() *round_types.Round
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

func (s *CarService) GenerateCar(element round_types.Element) *types.Car {
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
	car.Tank = r.Intn(max-min+1) + min
	car.Fuel = 10
	car.Health = 10

	return &car
}

func (s *CarService) СountSpeed(roundElement round_types.Element, car *types.Car) float64 {
	if car == nil {
		return 0
	}
	accelerationCoef := float64(car.Acceleration) / 100 * 25
	hendlingCoef := float64(car.Hendling) / 100 * 25
	brakesCoef := float64(car.Brakes) / 100 * 25
	strengthCoef := float64(car.Strength) / 100 * 25

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

	elementCoef := types.ElementEffects[car.Element][roundElement]

	return float64(accelerationCoef+hendlingCoef+brakesCoef+strengthCoef) * (float64(elementCoef) / 100)
}

func (s *CarService) countMiles(speed float64, duration int64, car *types.Car) float64 {
	return speed * (float64(duration) / 60 / 60)
}

func (s *CarService) countFuelWasting(speed float64, duration int64, roundElement round_types.Element, car *types.Car) float64 {

	wasting := speed * 26 / 10 * float64(car.Tank) / 100 / (float64(types.ElementEffectsFuel[car.Element][roundElement]) / 100)
	miles := s.countMiles(speed, duration, car)

	fmt.Printf("speed: %f, duration: %d, roundElement: %s, car: %+v\n", speed, duration, roundElement, car)
	fmt.Printf("wasting: %f, miles: %f\n", wasting, miles)

	wasted := miles / wasting

	return wasted
}

func (s *CarService) countHealthWasting(speed float64, duration int64, roundElement round_types.Element, car *types.Car) float64 {

	wasting := speed * 52 / 10 * float64(car.Strength) / 100 / (float64(types.ElementEffectsFuel[car.Element][roundElement]) / 100)

	miles := s.countMiles(speed, duration, car)

	wasted := miles / wasting

	return wasted
}

func (s *CarService) countDurationFromFuelWasting(speed float64, wasted float64, roundElement round_types.Element, car *types.Car) int64 {
	fmt.Printf("test speed: %f, wasted: %f, roundElement: %f, car: %+v\n", speed, wasted, float64(types.ElementEffectsFuel[car.Element][roundElement])/100, car)
	wasting := speed * 26 / 10 * float64(car.Tank) / 100 / (float64(types.ElementEffectsFuel[car.Element][roundElement]) / 100)
	miles := wasted * wasting
	fmt.Println(miles, speed)
	duration := int64(miles / speed * 60 * 60)
	return duration
}

func (s *CarService) countDurationFromHealthWasting(speed float64, wasted float64, roundElement round_types.Element, car *types.Car) int64 {
	wasting := speed * 52 / 10 * float64(car.Strength) / 100 / (float64(types.ElementEffectsFuel[car.Element][roundElement]) / 100)
	miles := wasted * wasting
	duration := int64(miles / speed * 60 * 60)
	return duration
}

func (s *CarService) GetAllCars(ctx context.Context) []types.Car {
	cars := []types.Car{}

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
	car := s.GenerateCar(element)
	car.IsMain = true

	if err := s.repo.BuyCar(ctx, car, userID); err != nil {
		return err
	}

	return nil
}

func (s *CarService) GetMainCar(ctx context.Context, userID int64) (*types.Car, error) {
	car, err := s.repo.GetMainCar(ctx, userID)
	if err != nil {
		return nil, err
	}

	if car != nil {
		round := s.roundService.GetRound()
		car.Speed = s.СountSpeed(round.Element, car)
	}
	return car, nil
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

func (s *CarService) GetRace(ctx context.Context, userID int64) (race *types.RaceState, err error) {
	round := s.roundService.GetRound()
	// TODO: check if fuel and health is enough
	return s.repo.GetRaceState(ctx, userID, round.ID)
}
func (s *CarService) StartRace(ctx context.Context, userID int64) (race *types.RaceState, err error) {
	round := s.roundService.GetRound()
	// TODO: check if fuel and health is enough
	return s.repo.StartRace(ctx, userID, round.ID)
}
func (s *CarService) EndRace(ctx context.Context, userID int64) (race *types.RaceState, err error) {
	round := s.roundService.GetRound()
	car, err := s.repo.GetMainCar(ctx, userID)
	if err != nil {
		return nil, err
	}

	state, err := s.repo.GetRaceState(ctx, userID, round.ID)
	if err != nil {
		return nil, err
	}

	speed := s.СountSpeed(round.Element, car)

	now := time.Now().Unix()
	raceTime := now - state.StartTime

	fuelWasted := s.countFuelWasting(speed, raceTime, round.Element, car)
	healthWasted := s.countHealthWasting(speed, raceTime, round.Element, car)

	fuelLeft := car.Fuel - fuelWasted
	healthLeft := car.Health - healthWasted

	fmt.Printf("fuelLeft: %f, healthLeft: %f\n", fuelLeft, healthLeft)

	// TODO: transfer to separate function
	// Check if car was broken during the way
	if fuelLeft <= 0 && healthLeft <= 0 {
		if healthLeft < fuelLeft {
			raceTime = s.countDurationFromHealthWasting(speed, healthLeft, round.Element, car)
		} else {
			raceTime = s.countDurationFromFuelWasting(speed, fuelLeft, round.Element, car)
		}
		fuelWasted = s.countFuelWasting(speed, raceTime, round.Element, car)
		healthWasted = s.countHealthWasting(speed, raceTime, round.Element, car)
	} else if fuelLeft <= 0 {
		raceTime = s.countDurationFromFuelWasting(speed, car.Fuel, round.Element, car)
		fuelWasted = s.countFuelWasting(speed, raceTime, round.Element, car)
		healthWasted = s.countHealthWasting(speed, raceTime, round.Element, car)
	} else if healthLeft <= 0 {
		raceTime = s.countDurationFromHealthWasting(speed, car.Health, round.Element, car)
		fuelWasted = s.countFuelWasting(speed, raceTime, round.Element, car)
		healthWasted = s.countHealthWasting(speed, raceTime, round.Element, car)
	}

	miles := s.countMiles(speed, raceTime, car)

	if err = s.repo.EndRace(ctx, userID, round.ID, miles, fuelWasted, healthWasted); err != nil {
		return nil, err
	}

	return s.repo.GetRaceState(ctx, userID, round.ID)

}

func (s *CarService) GetModulesOfUser(ctx context.Context, carID int64) ([]types.Module, error) {
	return s.repo.GetModulesOfUser(ctx, carID)
}

const (
	FuelCost   = 1
	HealthCost = 2
)

func (s *CarService) BuyFuel(ctx context.Context, userID int64) error {
	car, err := s.repo.GetMainCar(ctx, userID)
	if err != nil {
		return err
	}

	amount := 10 - int(car.Fuel)

	if err := s.repo.BuyFuel(ctx, userID, amount*FuelCost); err != nil {
		return err
	}

	return nil
}

func (s *CarService) BuyHealth(ctx context.Context, userID int64) error {
	car, err := s.repo.GetMainCar(ctx, userID)
	if err != nil {
		return err
	}

	amount := 10 - int(car.Health)

	if err := s.repo.BuyHealth(ctx, userID, amount*HealthCost); err != nil {
		return err
	}

	return nil
}
