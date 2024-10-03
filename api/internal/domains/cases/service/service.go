package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/cases/types"
)

// Errors
var (
	ErrNoCaseAvailable = errors.New("no case available")
	ErrCaseCannotOpen  = errors.New("case cannot open")
)

type repository interface {
	HasCase(userID int64, caseType string) (bool, error)
	OpenCase(userID int64, case_type *types.Case, reward *types.Reward) error

	GetCases(userID int64) ([]types.Case, error)
}

type userService interface {
}

type CaseService struct {
	repo        repository
	userService userService
}

func New(repo repository, userService userService) *CaseService {
	return &CaseService{
		repo:        repo,
		userService: userService,
	}
}

func (s *CaseService) OpenCase(userID int64, caseType string) (*types.Reward, error) {
	hasCase, err := s.repo.HasCase(userID, caseType)
	if err != nil {
		return nil, err
	}

	if !hasCase {
		return nil, ErrNoCaseAvailable
	}

	caseToOpen := types.Cases[caseType]

	reward, err := s.countRewards(caseType)
	if err != nil {
		return nil, err
	}

	if err = s.repo.OpenCase(userID, &caseToOpen, reward); err != nil {
		return nil, err
	}

	reward.Type = caseToOpen.RewardType
	return reward, nil
}

func (s *CaseService) countRewards(caseType string) (reward *types.Reward, err error) {
	caseToOpen := types.Cases[caseType]

	totalProbability := 0
	for _, reward := range caseToOpen.Rewards {
		totalProbability += reward.Probability
	}

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(totalProbability)

	probability := 0

	for _, reward := range caseToOpen.Rewards {
		probability += reward.Probability
		if randomNumber < probability {
			return &reward, nil
		}
	}
	return nil, ErrCaseCannotOpen
}

func (s *CaseService) GetCases(userID int64) ([]types.Case, error) {
	cases, err := s.repo.GetCases(userID)
	if err != nil {
		return nil, err
	}
	for i := range cases {
		cases[i] = types.Cases[cases[i].Type]
	}
	return cases, nil
}
