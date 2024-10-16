package service

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/round/types"
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

type newRoundSubscriber interface {
	AcceptNewRound(round *types.Round)
}

type RoundService struct {
	repo repository

	subscribers []newRoundSubscriber
}

func New(repo repository) *RoundService {
	return &RoundService{
		repo: repo,
	}
}

func (s *RoundService) RegisterSubscriber(subscriber newRoundSubscriber) {
	var mu sync.Mutex
	mu.Lock()
	s.subscribers = append(s.subscribers, subscriber)
	mu.Unlock()
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

	return round
}

func (s *RoundService) RunRounds() {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	// retriesNumber := 0
	// for {
	// 	if retriesNumber > 10 {
	// 		panic("couldn't start new round: error processing bets")
	// 	}
	// 	if err := s.battleService.GetNewBattles(); err != nil {
	// 		slog.Error("error processing bets: " + err.Error())
	// 		retriesNumber++
	// 		continue
	// 	}
	// 	break
	// }

	// go s.SetUpdateInterval()

	for {
		if time.Now().Unix() >= nextRoundStartTime {
			break
		}
	}
	s.StartNextRoundTimer()
}

func (s *RoundService) CurrentRound() *types.Round {
	roundID, endTime := s.countRound()
	// TODO: get element
	return &types.Round{
		ID:      roundID,
		EndTime: endTime,
	}
}

func (s *RoundService) StartRound() {
	time.Sleep(5 * time.Second)

	for _, sub := range s.subscribers {
		sub.AcceptNewRound(s.CurrentRound())
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
