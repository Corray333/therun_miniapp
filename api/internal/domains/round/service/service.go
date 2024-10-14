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
}

type battleService interface {
	GetBattles(round int, userID int64) ([]types.Battle, error)
	ProcessBets(roundID int) error
	GetNewBattles() error
	GetBattlesByID(ids []int) ([]types.Battle, error)

	UpdateBattles() error
}

type userService interface {
	GetUser(userID int64) (*user_types.User, error)
	ActivateUser(userID int64) error
}

type RoundService struct {
	repo          repository
	userService   userService
	battleService battleService
}

func New(repo repository, userService userService, battleService battleService) *RoundService {
	return &RoundService{
		repo:          repo,
		userService:   userService,
		battleService: battleService,
	}
}

func (s *RoundService) countRound() (roundID int, roundEndTime int64) {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	return currentRound, nextRoundStartTime
}

func (s *RoundService) GetRound(userID int64) *types.Round {
	roundID, roundEndTime := s.countRound()
	round := &types.Round{
		ID:      roundID,
		EndTime: roundEndTime,
	}

	battles, err := s.battleService.GetBattles(roundID, userID)
	if err != nil {
		slog.Error("error getting battles: " + err.Error())
		return round
	}

	round.Battles = battles
	return round
}

func (s *RoundService) RunRounds() {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	retriesNumber := 0
	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error processing bets")
		}
		if err := s.battleService.GetNewBattles(); err != nil {
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

func (s *RoundService) StartRound() {
	time.Sleep(5 * time.Second)

	retriesNumber := 0
	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error processing bets")
		}
		if err := s.battleService.UpdateBattles(); err != nil {
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
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	retriesNumber = 0
	for {
		if retriesNumber > 10 {
			panic("couldn't start new round: error getting new battles")
		}
		if err := s.battleService.GetNewBattles(); err != nil {
			slog.Error("error getting new battles: " + err.Error())
			retriesNumber++
			continue
		}
		return
	}
}

func (s *RoundService) StartNextRoundTimer() {
	// TODO: change to loop with time check
	slog.Info("Round starting...")
	go s.StartRound()
	time.AfterFunc(RoundDuration*time.Second, func() {
		s.StartNextRoundTimer()
	})
}

func (s *RoundService) SetUpdateInterval() {
	go s.battleService.GetNewBattles()
	time.AfterFunc(5*time.Minute, func() {
		s.SetUpdateInterval()
	})
}

func (s *RoundService) ProcessBets(roundID int) error {
	slog.Info("Processing bets...")
	return s.battleService.ProcessBets(roundID)
}
