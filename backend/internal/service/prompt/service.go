package prompt

import (
	"backend/internal/api/dto"
	"backend/internal/model"
	"backend/internal/repository/prompt"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

var (
	ErrPromptNotFound      = errors.New("prompt not found")
	ErrPromptAlreadyExists = errors.New("prompt already exists")
	ErrDatabaseErr         = errors.New("query error, please contact admin")
)

type IService interface {
	Create(ctx context.Context, p *model.Prompt) error
	Update(ctx context.Context, p *model.Prompt) error
	GetByID(ctx context.Context, id string) (*model.Prompt, error)
	GetByPath(ctx context.Context, path string) (*model.Prompt, error)
	List(ctx context.Context, userID string, offset, limit int) ([]*model.Prompt, int64, error)
	DeleteByID(ctx context.Context, id string) error
}

type Service struct {
	repo   *prompt.Repo
	logger *zap.Logger
}

func CreatePromptService(repo *prompt.Repo, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, req dto.CreatePromptDTO) (*model.Prompt, error) {
	p := &model.Prompt{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
		Username:  req.Username,
		Path:      req.Path,
		Category:  req.Category,
	}

	p, err := s.repo.Create(ctx, p)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, p *model.Prompt) error {
	old, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	if old == nil {
		return ErrPromptNotFound
	}

	p.Path = old.Path
	p.BaseModel.CreatedAt = old.CreatedAt
	p.BaseModel.UpdatedAt = time.Now()

	return s.repo.Update(ctx, p)
}

func (s *Service) GetByID(ctx context.Context, id string) (*model.Prompt, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return p, nil
}

func (s *Service) GetByPath(ctx context.Context, path string) (*model.Prompt, error) {
	p, err := s.repo.GetByPath(ctx, path)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	if p == nil {
		return nil, ErrPromptNotFound
	}
	return p, nil
}

func (s *Service) List(ctx context.Context, username string, offset, limit int) ([]*model.Prompt, int64, error) {
	p, err := s.repo.List(ctx, username, offset, limit)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, 0, ErrDatabaseErr
	}
	count, err := s.repo.CountByUser(ctx, username)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, 0, ErrDatabaseErr
	}
	return p, count, nil
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	old, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrPromptNotFound
	}
	if old == nil {
		return ErrPromptNotFound
	}
	err = s.repo.DeleteByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	return nil
}
