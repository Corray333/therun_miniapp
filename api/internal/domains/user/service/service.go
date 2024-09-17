package service

import (
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
	GetUser(userID int64) (*types.User, error)
	UpdateUser(user *types.User) error
	CreateUser(user *types.User) (err error)
	CheckIfRefCodeExists(refCode string) (bool, error)
	GetRefererID(refCode string) (int64, error)
	ListReferals(userID int64) ([]types.Referal, error)
	CountReferals(userID int64) (int, error)
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

func (s *UserService) AuthUser(initData, refCode string) (accessToken string, isNew bool, err error) {
	token, err := auth.CreateAccessToken(initData)
	if err != nil {
		return "", false, err
	}

	creds, err := auth.ExtractCredentials(token)
	if err != nil {
		return "", false, err
	}

	_, err = s.repo.GetUser(creds.ID)
	if err != nil && err != sql.ErrNoRows {
		return "", false, err
	}
	if err == nil {
		return token, false, nil
	}

	user := &types.User{
		ID:           creds.ID,
		Username:     creds.Username,
		PointBalance: 0,
		RaceBalance:  0,
		KeyBalance:   0,
		LastClaim:    time.Now().Unix(),
		MaxPoints:    MaxPointsDefault,
		FarmingTime:  FarmingTimeDefault,
	}

	avatar, err := s.external.GetAvatar(creds.ID)
	if err != nil {
		return "", false, err
	}

	if avatar != nil {
		filePath, err := s.fileManager.UploadImage(avatar, creds.Username)
		if err != nil {
			return "", false, err
		}
		user.Avatar = strings.TrimPrefix(filePath, "..")
	}

	if refCode != "" {
		referer, err := s.repo.GetRefererID(refCode)
		if err != nil && err != sql.ErrNoRows {
			return "", false, err
		}
		if referer != 0 {
			user.Referer = &referer
		}
	}

	numberOfRetries := 0
	for {
		user.RefCode, err = s.GenerateRefCode()
		if err != nil {
			return "", false, err
		}

		if err := s.repo.CreateUser(user); err == nil {
			return token, true, nil
		} else {
			slog.Error("error creating user: " + err.Error())
		}
		numberOfRetries++
		if numberOfRetries > MaxNumberOfRetries {
			return "", false, err
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

func (s *UserService) ListReferals(userID int64) ([]types.Referal, error) {
	refs, err := s.repo.ListReferals(userID)
	if err != nil {
		return nil, err
	}

	for i := range refs {
		refs[i].Avatar = os.Getenv("BASE_URL") + refs[i].Avatar
	}

	return refs, nil
}

const (
	Level1Referals = 3
	Level2Referals = 5
	Level3Referals = 10
	Level4Referals = 20
	Level5Referals = 50
	Level6Referals = 100
	Level7Referals = 200
	Level8Referals = 500
	Level9Referals = 1000
)

var levels = [9]int{3, 5, 10, 20, 50, 100, 200, 500, 1000}

func (s *UserService) CountReferals(userID int64) (count, level, nextLevelCount, previousLevelCount int, err error) {
	count, err = s.repo.CountReferals(userID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	level = 0
	nextLevelCount = levels[0]

	for i, l := range levels {
		if count >= l {
			level = i + 1
			if i == len(levels)-1 {
				nextLevelCount = 1000000
			} else {
				nextLevelCount = levels[i+1]
				previousLevelCount = levels[i]
			}
		}
	}

	return count, level + 1, nextLevelCount, previousLevelCount, nil
}
