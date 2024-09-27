package service

import (
	"log/slog"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/battle/types"
)

// Errors
var ()

const (
	BetAmount = 300 // Points
	BetPrize  = 1   // Keys
)

const (
	InitialRoundID        = 447
	InitialRoundStartTime = 1727283600
	RoundDuration         = 26 * 60 * 60 // 26 hours
)

type repository interface {
}

type external interface {
}

type BattleService struct {
	repo     repository
	external external
}

func New(repo repository, ext external) *BattleService {
	return &BattleService{
		repo:     repo,
		external: ext,
	}
}

func (s *BattleService) countRound() (roundID, roundEndTime int) {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	return currentRound, int(nextRoundStartTime)
}

func (s *BattleService) GetRound() *types.Round {
	roundID, roundEndTime := s.countRound()
	return &types.Round{
		ID:      roundID,
		EndTime: roundEndTime,
	}
}

func (s *BattleService) RunRounds() {
	elapsedTime := time.Now().Unix() - InitialRoundStartTime
	currentRound := InitialRoundID + int(elapsedTime/RoundDuration)
	currentRoundStartTime := InitialRoundStartTime + int64(currentRound-InitialRoundID)*RoundDuration

	nextRoundStartTime := currentRoundStartTime + RoundDuration

	for {
		if time.Now().Unix() >= nextRoundStartTime {
			break
		}
	}
	s.StartNextRoundTimer()
}

func (s *BattleService) StartRound() error {
	return nil
}

func (s *BattleService) StartNextRoundTimer() {
	time.AfterFunc(RoundDuration, func() {
		if err := s.StartRound(); err != nil {
			slog.Error("Failed to start next round" + err.Error())
			panic(err)
		}
		s.StartNextRoundTimer() // Повторяем для следующего раунда
	})
}
