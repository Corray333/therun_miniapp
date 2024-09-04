package service

import (
	"database/sql"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/Corray333/therun_miniapp/internal/domains/user/utils"
	"github.com/Corray333/therun_miniapp/pkg/server/auth"
)

const MaxPointsDefault = 200
const FarmTimeDefault = 7200
const MaxNumberOfRetries = 10

type repository interface {
	GetUser(userID int64) (*types.User, error)
	UpdateUser(user *types.User) error
	CreateUser(user *types.User) (err error)
	CheckIfRefCodeExists(refCode string) (bool, error)
	GetRefererID(refCode string) (int64, error)
	ListReferals(userID int64) ([]types.Referal, error)
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

func (s *UserService) AuthUser(initData, refCode string) (accessToken string, err error) {
	token, err := auth.CreateAccessToken(initData)
	if err != nil {
		return "", err
	}

	creds, err := auth.ExtractCredentials(token)
	if err != nil {
		return "", err
	}

	_, err = s.repo.GetUser(creds.ID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == nil {
		return token, nil
	}

	user := &types.User{
		ID:           creds.ID,
		Username:     creds.Username,
		PointBalance: 0,
		RaceBalance:  0,
		KeyBalance:   0,
		LastClaim:    time.Now().Unix(),
		MaxPoints:    MaxPointsDefault,
		FarmTime:     FarmTimeDefault,
	}

	avatar, err := s.external.GetAvatar(creds.ID)
	if err != nil {
		return "", err
	}

	if avatar != nil {
		filePath, err := s.fileManager.UploadImage(avatar, creds.Username)
		if err != nil {
			return "", err
		}
		user.Avatar = filePath
	}

	if refCode != "" {
		referer, err := s.repo.GetRefererID(refCode)
		if err != nil {
			return "", err
		}
		user.Referer = &referer
	}

	numberOfRetries := 0
	for {
		user.RefCode, err = s.GenerateRefCode()
		if err != nil {
			return "", err
		}

		if err := s.repo.CreateUser(user); err == nil {
			// TODO: maybe notify referer
			return token, nil
		}
		numberOfRetries++
		if numberOfRetries > MaxNumberOfRetries {
			return "", err
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
	return s.repo.ListReferals(userID)
}
