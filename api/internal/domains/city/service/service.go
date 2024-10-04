package service

import (
	"errors"
)

// Errors
var (
	ErrNoCaseAvailable = errors.New("no case available")
	ErrCaseCannotOpen  = errors.New("case cannot open")
)

type repository interface {
}

type CityService struct {
	repo repository
}

func New(repo repository) *CityService {
	return &CityService{
		repo: repo,
	}
}
