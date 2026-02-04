package favorites

import (
	"backend/internal/api/dto"
	"backend/internal/model"
	"backend/internal/repository/favorites"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrFavoriteNotFound = errors.New("favorite not found")
	ErrFavoriteExists   = errors.New("favorite already exists")
	ErrDatabaseErr      = errors.New("query error, please contact admin")
)

type IService interface {
	AddFavorite(ctx context.Context, userID int64, req dto.AddFavoriteDTO) (*model.Favorite, error)
	RemoveFavorite(ctx context.Context, userID int64, req dto.RemoveFavoriteDTO) error
	IsFavorite(ctx context.Context, userID int64, promptID string) (bool, error)
	ListFavorites(ctx context.Context, userID int64, offset, limit int) ([]*model.FavoriteWithPrompt, int64, error)
}

type Service struct {
	repo   *favorites.Repo
	logger *zap.Logger
}

func CreateFavoriteService(repo *favorites.Repo, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) AddFavorite(ctx context.Context, userID int64, req dto.AddFavoriteDTO) (*model.Favorite, error) {
	exists, err := s.repo.Exists(ctx, userID, req.PromptID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	if exists {
		return nil, ErrFavoriteExists
	}

	f := &model.Favorite{
		ID:       uuid.New().String(),
		UserID:   userID,
		PromptID: req.PromptID,
	}

	if err := s.repo.Create(ctx, f); err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return f, nil
}

func (s *Service) RemoveFavorite(ctx context.Context, userID int64, req dto.RemoveFavoriteDTO) error {
	err := s.repo.DeleteByUserAndPrompt(ctx, userID, req.PromptID)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	return nil
}

func (s *Service) IsFavorite(ctx context.Context, userID int64, promptID string) (bool, error) {
	return s.repo.Exists(ctx, userID, promptID)
}

func (s *Service) ListFavorites(ctx context.Context, userID int64, offset, limit int) ([]*model.FavoriteWithPrompt, int64, error) {
	list, err := s.repo.ListByUserWithPrompt(ctx, userID, offset, limit)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, 0, ErrDatabaseErr
	}
	count, err := s.repo.CountByUser(ctx, userID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, 0, ErrDatabaseErr
	}
	return list, count, nil
}
