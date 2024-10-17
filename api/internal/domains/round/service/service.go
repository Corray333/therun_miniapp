package service

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/round/types"
	"golang.org/x/exp/rand"
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
	SetRound(round *types.Round) error
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

func (s *RoundService) generateElement() types.Element {
	elements := types.Elements
	seed := uint64(time.Now().UnixNano())
	r := rand.New(rand.NewSource(seed))

	return elements[r.Intn(len(elements))]
}

func (s *RoundService) RunRounds() {
	roundID, endTime := s.countRound()

	s.repo.SetRound(&types.Round{
		ID:      roundID,
		EndTime: endTime,
		Element: s.generateElement(),
	})

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
		if time.Now().Unix() >= endTime {
			break
		}
	}
	s.StartNextRoundTimer()
}

func (s *RoundService) GetRound() *types.Round {

	// TODO: get element from db
	roundID, roundEndTime := s.countRound()
	round := &types.Round{
		ID:      roundID,
		EndTime: roundEndTime,
		Element: "desert",
	}

	return round
}

func (s *RoundService) StartRound() {
	time.Sleep(5 * time.Second)

	go func() {
		s.repo.SetRound(&types.Round{
			ID:      s.GetRound().ID,
			EndTime: s.GetRound().EndTime,
			Element: s.generateElement(),
		})
	}()

	for _, sub := range s.subscribers {
		sub.AcceptNewRound(s.GetRound())
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
