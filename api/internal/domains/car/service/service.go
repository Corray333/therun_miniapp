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

func (s *CarService) СountSpeed(car *types.Car) float64 {
	accelerationCoef := float64(car.Acceleration) / 100 * 25
	hendlingCoef := float64(car.Hendling) / 100 * 25
	brakesCoef := float64(car.Brakes) / 100 * 25
	strengthCoef := float64(car.Strength) / 100 * 25

	fmt.Println(accelerationCoef, hendlingCoef, brakesCoef, strengthCoef)

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

	round := s.roundService.GetRound()
	elementCoef := types.ElementEffects[car.Element][round.Element]

	fmt.Println(car.Element, round.Element)
	fmt.Println("Elem coef: ", elementCoef)

	return float64(accelerationCoef+hendlingCoef+brakesCoef+strengthCoef) * (float64(elementCoef) / 100)
}

func (s *CarService) countMiles(duration int64, car *types.Car) float64 {
	speed := s.СountSpeed(car)
	return speed * float64(duration) / 60 / 60
}

func (s *CarService) countFuelWasting(duration int64, car *types.Car) float64 {
	round := s.roundService.GetRound()
	speed := s.СountSpeed(car)

	fmt.Println(speed)
	fmt.Println(car)
	wasting := speed * 26 / 10 * float64(car.Tank) / 100 / float64(types.ElementEffectsFuel[car.Element][round.Element]) / 100
	fmt.Println("Wasting: ", wasting)
	miles := s.countMiles(duration, car)

	wasted := miles / wasting

	return wasted
}

func (s *CarService) countHealthWasting(duration int64, car *types.Car) float64 {
	speed := s.СountSpeed(car)

	wasting := speed * 52 / 10 * float64(car.Strength) / 100

	miles := s.countMiles(duration, car)

	wasted := miles / wasting

	return wasted
}

// TODO: rewrite using round info
func (s *CarService) countDurationFromFuelWasting(wasting float64, car *types.Car) int64 {
	// round := s.roundService.CurrentRound()
	speed := s.СountSpeed(car)

	miles := wasting * speed * 26 / 10 * 100 / float64(car.Fuel)

	return int64(miles / speed * 60 * 60)
}

// TODO: rewrite using round info
func (s *CarService) countDurationFromHealthWasting(wasting float64, car *types.Car) int64 {
	speed := s.СountSpeed(car)

	miles := wasting * speed * 52 / 10 * 100 / float64(car.Strength)

	return int64(miles / speed * 60 * 60)
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

	now := time.Now().Unix()
	raceTime := now - state.StartTime

	fuelWasted := s.countFuelWasting(raceTime, car)
	healthWasted := s.countHealthWasting(raceTime, car)

	fmt.Printf("fuelWasted: %f, healthWasted: %f, timeSpent: %d\n", fuelWasted, healthWasted, raceTime)

	fuelLeft := car.Fuel - fuelWasted
	healthLeft := car.Health - healthWasted

	// TODO: transfer to separate function
	// Check if car was broken during the way
	if fuelLeft <= 0 && healthLeft <= 0 {
		if healthLeft < fuelLeft {
			raceTime = s.countDurationFromHealthWasting(healthLeft, car)
		} else {
			raceTime = s.countDurationFromFuelWasting(fuelLeft, car)
		}
		fuelWasted = s.countFuelWasting(raceTime, car)
		healthWasted = s.countHealthWasting(raceTime, car)
	} else if fuelLeft <= 0 {
		raceTime = s.countDurationFromFuelWasting(fuelLeft, car)
		fuelWasted = s.countFuelWasting(raceTime, car)
		healthWasted = s.countHealthWasting(raceTime, car)
	} else if healthLeft <= 0 {
		raceTime = s.countDurationFromHealthWasting(healthLeft, car)
		fuelWasted = s.countFuelWasting(raceTime, car)
		healthWasted = s.countHealthWasting(raceTime, car)
	}

	miles := s.countMiles(raceTime, car)

	if err = s.repo.EndRace(ctx, userID, round.ID, miles, fuelWasted, healthWasted); err != nil {
		return nil, err
	}

	return s.repo.GetRaceState(ctx, userID, round.ID)

}
