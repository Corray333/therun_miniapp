package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
	round_service "github.com/Corray333/therun_miniapp/internal/domains/round/service"
	round_types "github.com/Corray333/therun_miniapp/internal/domains/round/types"
	user_types "github.com/Corray333/therun_miniapp/internal/domains/user/types"
)

// Errors
var (
	ErrTooLate         = errors.New("too late to make a bet")
	ErrNotEnoughPoints = errors.New("not enough points")
)

const (
	BetAmount = 300 // Points
	BetPrize  = 1   // Keys
)

const (
	InitialRoundID        = 447
	InitialRoundStartTime = 1727283600
	RoundDuration         = 26 * 60 * 60 // 26 hours
	BetsPeriodDuration    = 12 * 60 * 60 // 12 hours
)

type repository interface {
	SetBattles(battles []types.Battle) error
	GetBattles(round int, userID int64) ([]types.Battle, error)

	GetBattle(id int) (*types.Battle, error)
	MakeBet(userID int64, battleID, pick int) error

	ProcessBets(roundID int, keyReward int) error
}

type external interface {
	GetNewBattles() ([]types.Battle, error)
	GetBattlesByID(ids []int) ([]types.Battle, error)
}

type userService interface {
	GetUser(userID int64) (*user_types.User, error)
	ActivateUser(userID int64) error
}

type roundService interface {
	GetRound() (*round_types.Round, error)
}

type BattleService struct {
	repo         repository
	external     external
	userService  userService
	roundService roundService
}

func New(repo repository, ext external, userService userService, roundService *round_service.RoundService) *BattleService {
	return &BattleService{
		repo:         repo,
		external:     ext,
		userService:  userService,
		roundService: roundService,
	}
}

func (s *BattleService) AcceptNewRound(round *round_types.Round) {
	retriesNumber := 0
	for {
		if retriesNumber > 10 {
			slog.Error("couldn't start new round: error updating battles")
			break
		}
		if err := s.UpdateBattles(); err != nil {
			slog.Error("error processing bets: " + err.Error())
			retriesNumber++
			continue
		}
		break
	}
	retriesNumber = 0

	for {
		if retriesNumber > 10 {
			slog.Error("couldn't start new round: error processing bets")
			break
		}
		if err := s.ProcessBets(round.ID - 1); err != nil {
			slog.Error("error processing bets: " + err.Error())
			retriesNumber++
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	retriesNumber = 0
	for {
		if retriesNumber > 10 {
			slog.Error("couldn't start new round: error getting new battles")
			break
		}
		if err := s.GetNewBattles(); err != nil {
			slog.Error("error getting new battles: " + err.Error())
			retriesNumber++
			continue
		}
		return
	}
}

func (s *BattleService) GetNewBattles() error {
	slog.Info("Getting new battles...")
	battles, err := s.external.GetNewBattles()
	if err != nil {
		return err
	}

	if err := s.repo.SetBattles(battles); err != nil {
		return err
	}

	return nil
}

func (s *BattleService) MakeBet(userID int64, battleID, pick int) error {
	user, err := s.userService.GetUser(userID)
	if err != nil {
		return err
	}

	if user.PointBalance < BetAmount {
		return ErrNotEnoughPoints
	}

	if !user.IsActivated {
		if err := s.userService.ActivateUser(userID); err != nil {
			return err
		}
	}

	round, err := s.roundService.GetRound()
	if err != nil {
		return err
	}
	battle, err := s.repo.GetBattle(battleID)
	if err != nil {
		return err
	}

	if battle.RoundID != round.ID || time.Now().Unix() >= round.EndTime-BetsPeriodDuration {
		return ErrTooLate
	}

	return s.repo.MakeBet(userID, battleID, pick)
}

func (s *BattleService) UpdateBattles() error {
	slog.Info("Updating battles...")
	round, err := s.roundService.GetRound()
	if err != nil {
		return err
	}
	battles, err := s.repo.GetBattles(round.ID, 0)
	if err != nil {
		return err
	}

	ids := make([]int, len(battles))
	for i, battle := range battles {
		ids[i] = battle.ID
	}

	battles, err = s.external.GetBattlesByID(ids)
	if err != nil {
		return err
	}

	if err := s.repo.SetBattles(battles); err != nil {
		return err
	}

	return nil
}

func (s *BattleService) SetUpdateInterval() {
	go s.GetNewBattles()
	time.AfterFunc(5*time.Minute, func() {
		s.SetUpdateInterval()
	})
}

func (s *BattleService) ProcessBets(roundID int) error {
	slog.Info("Processing bets...")
	return s.repo.ProcessBets(roundID, BetPrize)
}

func (s *BattleService) GetBattles(userID int64) ([]types.Battle, error) {
	round, err := s.roundService.GetRound()
	if err != nil {
		return nil, err
	}
	fmt.Println(round)
	return s.repo.GetBattles(round.ID, userID)
}

func (s *BattleService) GetBattlesByID(ids []int) ([]types.Battle, error) {
	return s.external.GetBattlesByID(ids)
}
