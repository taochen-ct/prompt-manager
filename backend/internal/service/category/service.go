package category

import (
	"backend/internal/api/dto"
	"backend/internal/model"
	"backend/internal/repository/category"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrDatabaseErr           = errors.New("query error, please contact admin")
)

type IService interface {
	Create(ctx context.Context, req dto.CreateCategoryDTO) (*model.Category, error)
	Update(ctx context.Context, c *model.Category) error
	GetByID(ctx context.Context, id string) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
	DeleteByID(ctx context.Context, id string) error
}

type Service struct {
	repo   *category.Repo
	logger *zap.Logger
}

func CreateCategoryService(repo *category.Repo, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, req dto.CreateCategoryDTO) (*model.Category, error) {
	c := &model.Category{
		ID:    req.ID,
		Title: req.Title,
		Icon:  req.Icon,
		URL:   req.URL,
	}

	c, err := s.repo.Create(ctx, c)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return c, nil
}

func (s *Service) Update(ctx context.Context, c *model.Category) error {
	old, err := s.repo.GetByID(ctx, c.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	if old == nil {
		return ErrCategoryNotFound
	}

	c.BaseModel.CreatedAt = old.CreatedAt
	c.BaseModel.UpdatedAt = time.Now()

	return s.repo.Update(ctx, c)
}

func (s *Service) GetByID(ctx context.Context, id string) (*model.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return c, nil
}

func (s *Service) List(ctx context.Context) ([]*model.Category, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return list, nil
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	if c == nil {
		return ErrCategoryNotFound
	}
	err = s.repo.DeleteByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	return nil
}
