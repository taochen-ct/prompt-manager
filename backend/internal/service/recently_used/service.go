package recently_used

import (
	"backend/internal/api/dto"
	"backend/internal/model"
	"backend/internal/repository/recently_used"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const DefaultKeepCount = 50

var (
	ErrDatabaseErr = errors.New("query error, please contact admin")
)

type IService interface {
	RecordUsage(ctx context.Context, userID int64, req dto.RecordUsageDTO) (*model.RecentlyUsed, error)
	RemoveRecord(ctx context.Context, userID int64, req dto.RemoveRecentlyUsedDTO) error
	ListRecentlyUsed(ctx context.Context, userID int64, offset, limit int) ([]*model.RecentlyUsedWithPrompt, int64, error)
	CleanOldRecords(ctx context.Context, userID int64, keepCount int) error
}

type Service struct {
	repo   *recently_used.Repo
	logger *zap.Logger
}

func CreateRecentlyUsedService(repo *recently_used.Repo, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) RecordUsage(ctx context.Context, userID int64, req dto.RecordUsageDTO) (*model.RecentlyUsed, error) {
	rec := &model.RecentlyUsed{
		ID:       uuid.New().String(),
		UserID:   userID,
		PromptID: req.PromptID,
	}

	rec, err := s.repo.CreateOrUpdate(ctx, rec)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}

	if err := s.repo.DeleteOldByUser(ctx, userID, DefaultKeepCount); err != nil {
		s.logger.Warn("failed to clean old records", zap.Error(err))
	}

	return rec, nil
}

func (s *Service) RemoveRecord(ctx context.Context, userID int64, req dto.RemoveRecentlyUsedDTO) error {
	err := s.repo.DeleteByUserAndPrompt(ctx, userID, req.PromptID)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	return nil
}

func (s *Service) ListRecentlyUsed(ctx context.Context, userID int64, offset, limit int) ([]*model.RecentlyUsedWithPrompt, int64, error) {
	if limit <= 0 {
		limit = 10
	}
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

func (s *Service) CleanOldRecords(ctx context.Context, userID int64, keepCount int) error {
	if keepCount <= 0 {
		keepCount = DefaultKeepCount
	}
	return s.repo.DeleteOldByUser(ctx, userID, keepCount)
}
