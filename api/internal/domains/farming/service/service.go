package service

import (
	"errors"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
)

var (
	ErrClaimTooEarly = errors.New("claim too early")
)

type repository interface {
}

type userService interface {
	GetUser(userID int64) (*types.User, error)
	UpdateUser(user *types.User) error
}

type FarmingService struct {
	repo        repository
	userService userService
}

func New(repo repository, userService userService) *FarmingService {
	return &FarmingService{
		repo:        repo,
		userService: userService,
	}
}

func (s *FarmingService) ClaimTokens(userID int64) (pointsGot, pointBalance, farmingTime, maxPoints int, lastClaim int64, err error) {
	user, err := s.userService.GetUser(userID)
	if err != nil {
		return
	}

	if user.LastClaim+int64(user.FarmingTime) > time.Now().Unix() {
		return 0, user.PointBalance, user.FarmingTime, user.MaxPoints, user.LastClaim, ErrClaimTooEarly
	}

	user.PointBalance += user.MaxPoints
	user.LastClaim = time.Now().Unix()

	if err = s.userService.UpdateUser(user); err != nil {
		return
	}

	return user.MaxPoints, user.PointBalance, user.FarmingTime, user.MaxPoints, user.LastClaim, nil
}

func (s *FarmingService) GetUser(userID int64) (user *types.User, err error) {
	return s.userService.GetUser(userID)
}
