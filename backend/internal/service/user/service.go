package user

import (
	"backend/internal/repository/user"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"

	"backend/internal/api/dto"
	"backend/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("username already exists")
	ErrService      = errors.New("else error, please contact admin")
)

type IService interface {
	Create(ctx context.Context, req dto.CreateUserDTO) (*model.User, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]model.User, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserDTO) error
	Delete(ctx context.Context, id int64) error
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

	userData := &model.User{
		Username:   req.Username,
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
