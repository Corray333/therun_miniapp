package service

import (
	"math/rand"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/car/types"
	round_types "github.com/Corray333/therun_miniapp/internal/domains/round/types"
)

type repository interface {
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

func (s *CarService) countMiles(duration int, car *types.Car) float64 {
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

	return float64(accelerationCoef+hendlingCoef+brakesCoef+strengthCoef) * (float64(elementCoef) / 100) * float64(duration)
}
