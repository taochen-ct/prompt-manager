package user

import (
	"backend/internal/repository/user"
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/api/dto"
	"backend/internal/model"
	"backend/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("username already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrService         = errors.New("else error, please contact admin")
	ErrTokenExpired    = errors.New("token has expired")
	ErrInvalidToken    = errors.New("invalid token")
)

type LoginResponse struct {
	Token    string      `json:"token"`
	User     *model.User `json:"user"`
	ExpireAt time.Time   `json:"expireAt"`
}

type IService interface {
	Create(ctx context.Context, req dto.CreateUserDTO) (*model.User, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]*model.User, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserDTO) error
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, req dto.LoginDTO, secret string, expireHour int) (*LoginResponse, error)
	Logout(ctx context.Context, username string) error
}

type Service struct {
	repo   *user.Repo
	logger *zap.Logger
}

func CreateUserService(repo *user.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req dto.CreateUserDTO) (*model.User, error) {
	// 用户名唯一性检查
	if u, _ := s.repo.GetByUsername(ctx, req.Username); u != nil {
		return nil, ErrUserExists
	}

	// 密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrService
	}

	userData := &model.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Nickname:   req.Nickname,
		Department: req.Department,
	}

	if err := s.repo.Create(ctx, userData); err != nil {
		s.logger.Error(err.Error())
		return nil, ErrService
	}

	return userData, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*model.User, error) {
	userData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, ErrService
	}
	return userData, nil
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]*model.User, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *Service) Update(ctx context.Context, id int64, req dto.UpdateUserDTO) error {
	userData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	userData.Nickname = req.Nickname
	userData.Department = req.Department

	return s.repo.Update(ctx, userData)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *Service) Login(ctx context.Context, req dto.LoginDTO, secret string, expireHour int) (*LoginResponse, error) {
	userData, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		s.logger.Error(err.Error())
		return nil, ErrService
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// 生成 JWT
	expireDuration := time.Duration(expireHour) * time.Hour
	token, err := jwt.GenerateToken(userData.ID, userData.Username, secret, expireDuration)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrService
	}

	return &LoginResponse{
		Token:    token,
		User:     userData,
		ExpireAt: time.Now().Add(expireDuration),
	}, nil
}

func (s *Service) Logout(ctx context.Context, username string) error {
	return nil
}
