package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
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

type BattleService struct {
	repo        repository
	external    external
	userService userService
}

func New(repo repository, ext external, userService userService) *BattleService {
	return &BattleService{
		repo:        repo,
		external:    ext,
		userService: userService,
	}
}

func (s *BattleService) countRound() (roundID int, roundEndTime int64) {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	return currentRound, nextRoundStartTime
}

func (s *BattleService) GetRound(userID int64) *types.Round {
	roundID, roundEndTime := s.countRound()
	round := &types.Round{
		ID:      roundID,
		EndTime: roundEndTime,
	}

	battles, err := s.repo.GetBattles(roundID, userID)
	if err != nil {
		slog.Error("error getting battles: " + err.Error())
		return round
	}

	round.Battles = battles
	return round
}

func (s *BattleService) RunRounds() {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	retriesNumber := 0
	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error processing bets")
		}
		if err := s.GetNewBattles(); err != nil {
			slog.Error("error processing bets: " + err.Error())
			retriesNumber++
			continue
		}
		break
	}
	go s.SetUpdateInterval()

	for {
		if time.Now().Unix() >= nextRoundStartTime {
			break
		}
	}
	s.StartNextRoundTimer()
}

func (s *BattleService) StartRound() {
	time.Sleep(5 * time.Second)

	retriesNumber := 0
	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error processing bets")
		}
		if err := s.UpdateBattles(); err != nil {
			slog.Error("error processing bets: " + err.Error())
			retriesNumber++
			continue
		}
		break
	}
	retriesNumber = 0

	round, _ := s.countRound()

	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error processing bets")
		}
		if err := s.ProcessBets(round - 1); err != nil {
			slog.Error("error processing bets: " + err.Error())
			retriesNumber++
			continue
		}
		break
	}

	retriesNumber = 0
	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error getting new battles")
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

	roundID, roundEndTime := s.countRound()
	battle, err := s.repo.GetBattle(battleID)
	if err != nil {
		return err
	}

	if battle.RoundID != roundID || time.Now().Unix() >= roundEndTime-BetsPeriodDuration {
		return ErrTooLate
	}

	return s.repo.MakeBet(userID, battleID, pick)
}

func (s *BattleService) UpdateBattles() error {
	slog.Info("Updating battles...")
	round, _ := s.countRound()
	battles, err := s.repo.GetBattles(round, 0)
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

func (s *BattleService) StartNextRoundTimer() {
	// TODO: change to loop with time check
	slog.Info("Round starting...")
	go s.StartRound()
	time.AfterFunc(RoundDuration*time.Second, func() {
		s.StartNextRoundTimer()
	})
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
