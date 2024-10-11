package service

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/domains/user/utils"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
)

const MaxPointsDefault = 200
const FarmingTimeDefault = 7200
const MaxNumberOfRetries = 10

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetUser(userID int64) (*types.User, error)
	UpdateUser(user *types.User) error
	CreateUser(user *types.User) (err error)
	CheckIfRefCodeExists(refCode string) (bool, error)
	GetRefererID(refCode string) (int64, error)

	ListActivatedReferals(userID int64) ([]types.Referal, error)
	ListNotActivatedReferals(userID int64) ([]types.Referal, error)
	CountReferals(userID int64) (refsActivated, refsFrozenTotal, refsFrozen, refsPremiumFrozen, refsPremiumActivatedNotClaimed, refsActivatedNotClaimed int, err error)
	ClaimRefs(ctx context.Context, userID int64) (refsPremiumActivatedNotClaimed, refsActivatedNotClaimed int, err error)

	SetPremium(userID int64, isPremium bool) error

	ActivateUser(userID int64) error

	ChangeBalances(ctx context.Context, userID int64, changes []types.BalanceChange) error
}
type external interface {
	GetAvatar(userID int64) ([]byte, error)
}

type fileManager interface {
	UploadImage(file []byte, name string) (string, error)
}

type UserService struct {
	repo        repository
	external    external
	fileManager fileManager
}

func New(repo repository, external external, fileManager fileManager) *UserService {
	return &UserService{
		repo:        repo,
		external:    external,
		fileManager: fileManager,
	}
}

func (s *UserService) GetUser(userID int64) (*types.User, error) {
	return s.repo.GetUser(userID)
}

func (s *UserService) AuthUser(initData, refCode string) (accessToken string, isNew bool, isPremium bool, err error) {
	token, isPremium, err := auth.CreateAccessToken(initData)
	if err != nil {
		return "", false, false, err
	}

	creds, err := auth.ExtractCredentials(token)
	if err != nil {
		return "", false, false, err
	}

	_, err = s.repo.GetUser(creds.ID)
	if err != nil && err != sql.ErrNoRows {
		return "", false, false, err
	}
	if err == nil {
		return token, false, isPremium, nil
	}

	if creds.ID == 6202406149 {
		isPremium = true
	}

	pointBalance := 0
	var refererID *int64

	if refCode != "" {
		referer, err := s.repo.GetRefererID(refCode)
		if err != nil && err != sql.ErrNoRows {
			return "", false, false, err
		}
		if referer != 0 {
			refererID = &referer

			if isPremium {
				pointBalance = RefRewardPremium
			} else {
				pointBalance = RefReward
			}
		}
	}

	user := &types.User{
		ID:            creds.ID,
		Username:      creds.Username,
		PointBalance:  pointBalance,
		RaceBalance:   0,
		RedKeyBalance: 0,
		LastClaim:     time.Now().Unix(),
		MaxPoints:     MaxPointsDefault,
		FarmingTime:   FarmingTimeDefault,
		IsPremium:     isPremium,
		IsActivated:   isPremium,
		Referer:       refererID,
	}

	avatar, err := s.external.GetAvatar(creds.ID)
	if err != nil {
		return "", false, false, err
	}

	if avatar != nil {
		filePath, err := s.fileManager.UploadImage(avatar, creds.Username)
		if err != nil {
			return "", false, false, err
		}
		user.Avatar = strings.TrimPrefix(filePath, "..")
	}

	numberOfRetries := 0
	for {
		user.RefCode, err = s.GenerateRefCode()
		if err != nil {
			return "", false, false, err
		}

		if err := s.repo.CreateUser(user); err == nil {
			return token, true, isPremium, nil
		} else {
			slog.Error("error creating user: " + err.Error())
		}
		numberOfRetries++
		if numberOfRetries > MaxNumberOfRetries {
			return "", false, false, err
		}
	}
}

func (s *UserService) GenerateRefCode() (string, error) {
	numberOfRetries := 0
	for {
		code := utils.GenerateRefreshToken()
		exists, err := s.repo.CheckIfRefCodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
		numberOfRetries++
		if numberOfRetries > MaxNumberOfRetries {
			return "", err
		}
	}
}

func (s *UserService) UpdateUser(user *types.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) ListActivatedReferals(userID int64) ([]types.Referal, error) {
	refs, err := s.repo.ListActivatedReferals(userID)
	if err != nil {
		return nil, err
	}

	for i := range refs {
		if refs[i].Avatar != "" {
			refs[i].Avatar = os.Getenv("BASE_URL") + refs[i].Avatar
		}
	}

	return refs, nil
}

func (s *UserService) ListNotActivatedReferals(userID int64) ([]types.Referal, error) {
	refs, err := s.repo.ListNotActivatedReferals(userID)
	if err != nil {
		return nil, err
	}

	for i := range refs {
		if refs[i].Avatar != "" {
			refs[i].Avatar = os.Getenv("BASE_URL") + refs[i].Avatar
		}
	}

	return refs, nil
}

const (
	RefReward        = 1000
	RefRewardPremium = 3000
)

func (s *UserService) CountReferals(userID int64) (refsActivated, refsFrozen int, rewardsFrozenTotal, rewardsAvailibleTotal []types.BalanceChange, err error) {
	refsActivated, refsFrozenTotal, refsFrozen, refsPremiumFrozen, refsPremiumActivatedNotClaimed, refsActivatedNotClaimed, err := s.repo.CountReferals(userID)
	if err != nil {
		return 0, 0, nil, nil, err
	}

	// var refRewardFrozen, refRewardFrozenPremium, refRewardAvailible, refRewardAvailiblePremium = []types.BalanceChange{}, []types.BalanceChange{}, []types.BalanceChange{}, []types.BalanceChange{}
	refRewardFrozen := append([]types.BalanceChange{}, types.RefReward...)
	refRewardFrozenPremium := append([]types.BalanceChange{}, types.RefRewardPremium...)
	refRewardAvailible := append([]types.BalanceChange{}, types.RefReward...)
	refRewardAvailiblePremium := append([]types.BalanceChange{}, types.RefRewardPremium...)

	if refsFrozen == 0 {
		refRewardFrozen = nil
	}
	if refsPremiumFrozen == 0 {
		refRewardFrozenPremium = nil
	}
	if refsActivatedNotClaimed == 0 {
		refRewardAvailible = nil
	}
	if refsPremiumActivatedNotClaimed == 0 {
		refRewardAvailiblePremium = nil
	}

	for i := 0; i < len(refRewardFrozen); i++ {
		refRewardFrozen[i].Amount *= refsFrozen
	}

	for i := 0; i < len(refRewardFrozenPremium); i++ {
		refRewardFrozenPremium[i].Amount *= refsPremiumFrozen
	}

	for i := 0; i < len(refRewardAvailible); i++ {
		refRewardAvailible[i].Amount *= refsActivatedNotClaimed
	}

	for i := 0; i < len(refRewardAvailiblePremium); i++ {
		refRewardAvailiblePremium[i].Amount *= refsPremiumActivatedNotClaimed
	}

	rewardsFrozenMap := map[types.Currency]int{}
	for i := 0; i < len(refRewardFrozen); i++ {
		rewardsFrozenMap[refRewardFrozen[i].Currency] += refRewardFrozen[i].Amount
	}

	for i := 0; i < len(refRewardFrozenPremium); i++ {
		rewardsFrozenMap[refRewardFrozenPremium[i].Currency] += refRewardFrozenPremium[i].Amount
	}

	rewardsAvailibleMap := map[types.Currency]int{}
	for i := 0; i < len(refRewardAvailible); i++ {
		rewardsAvailibleMap[refRewardAvailible[i].Currency] += refRewardAvailible[i].Amount
	}

	for i := 0; i < len(refRewardAvailiblePremium); i++ {
		rewardsAvailibleMap[refRewardAvailiblePremium[i].Currency] += refRewardAvailiblePremium[i].Amount
	}

	rewardsFrozenTotal = []types.BalanceChange{}
	for currency, amount := range rewardsFrozenMap {
		rewardsFrozenTotal = append(rewardsFrozenTotal, types.BalanceChange{
			Currency: currency,
			Amount:   amount,
		})
	}

	rewardsAvailibleTotal = []types.BalanceChange{}
	for currency, amount := range rewardsAvailibleMap {
		rewardsAvailibleTotal = append(rewardsAvailibleTotal, types.BalanceChange{
			Currency: currency,
			Amount:   amount,
		})
	}

	return refsActivated, refsFrozenTotal, rewardsFrozenTotal, rewardsAvailibleTotal, nil

}

func (s *UserService) ClaimRefs(ctx context.Context, userID int64) (err error) {
	// rewardsGot = refsActivatedNotClaimed*service.RefReward + refsPremiumActivatedNotClaimed*service.RefRewardPremium
	ctx, err = s.repo.Begin(ctx)
	if err != nil {
		return err
	}

	refsPremiumActivatedNotClaimed, refsActivatedNotClaimed, err := s.repo.ClaimRefs(ctx, userID)
	if err != nil {
		return err
	}

	refsReward := []types.BalanceChange{}
	copy(refsReward, types.RefReward)

	refsRewardPremium := []types.BalanceChange{}
	copy(refsRewardPremium, types.RefRewardPremium)

	for i := 0; i < len(refsReward); i++ {
		refsReward[i].Amount *= refsActivatedNotClaimed
	}

	for i := 0; i < len(refsRewardPremium); i++ {
		refsRewardPremium[i].Amount *= refsPremiumActivatedNotClaimed
	}

	totalReward := append(refsReward, refsRewardPremium...)

	if err := s.repo.ChangeBalances(ctx, userID, totalReward); err != nil {
		return err
	}

	return s.repo.Commit(ctx)
}

const DayTime = 86400

var dayliCheckRewards = map[int]int{
	1: 100,
	2: 200,
	3: 300,
	4: 400,
	5: 500,
	6: 800,
	7: 1000,
}

func (s *UserService) DailyCheck(userID int64) (dailyCheckStreak int, dailyCheckLast int64, err error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now().UTC()
	lastCheckTime := time.Unix(user.DailyCheckLast, 0).UTC()

	// Check if last daily check was until 00:00 UTC today
	yesterday := now.AddDate(0, 0, -1)

	// Compare the date part of the given time with yesterday's date
	if lastCheckTime.Year() == yesterday.Year() && lastCheckTime.YearDay() == yesterday.YearDay() {
		user.DailyCheckStreak++
		if user.DailyCheckStreak > 7 {
			user.PointBalance += dayliCheckRewards[7]
		} else {
			user.PointBalance += dayliCheckRewards[user.DailyCheckStreak]
		}
	} else if !(lastCheckTime.Year() == yesterday.Year() && lastCheckTime.YearDay() == now.YearDay()) {
		user.DailyCheckStreak = 1
		user.PointBalance += dayliCheckRewards[1]
	}

	user.DailyCheckLast = now.Unix()

	if err := s.repo.UpdateUser(user); err != nil {
		return 0, 0, err
	}

	return user.DailyCheckStreak, user.DailyCheckLast, nil
}

func (s *UserService) SetPremium(userID int64, isPremium bool) error {
	return s.repo.SetPremium(userID, isPremium)
}

func (s *UserService) ActivateUser(userID int64) error {
	return s.repo.ActivateUser(userID)
}
