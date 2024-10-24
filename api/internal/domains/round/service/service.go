package service

import (
	"database/sql"
	"errors"
	"log/slog"
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
	GetRoundElement(roundID int) (types.Element, error)
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
	s.subscribers = append(s.subscribers, subscriber)
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
	_, err := s.repo.GetRoundElement(roundID)
	if err == sql.ErrNoRows {
		if err := s.repo.SetRound(&types.Round{
			ID:      roundID,
			EndTime: endTime,
			Element: s.generateElement(),
		}); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	for {
		if time.Now().Unix() >= endTime {
			break
		}
	}
	s.StartNextRoundTimer()
}

func (s *RoundService) GetRound() (*types.Round, error) {
	// TODO: get element from db
	roundID, roundEndTime := s.countRound()

	element, err := s.repo.GetRoundElement(roundID)
	if err != nil {
		return nil, err
	}

	round := &types.Round{
		ID:      roundID,
		EndTime: roundEndTime,
		Element: element,
	}

	return round, nil
}

func (s *RoundService) StartRound() {
	time.Sleep(time.Second)

	roundID, endTime := s.countRound()
	round := &types.Round{
		ID:      roundID,
		EndTime: endTime,
		Element: s.generateElement(),
	}

	if err := s.repo.SetRound(round); err != nil {
		slog.Error("error while setting new round: " + err.Error())
	}

	for _, sub := range s.subscribers {
		sub.AcceptNewRound(round)
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
