package version

import (
	"backend/internal/api/dto"
	"backend/internal/model"
	"backend/internal/repository/prompt"
	"backend/internal/repository/version"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrVersionNotFound      = errors.New("version not found")
	ErrPromptNotFound       = errors.New("prompt not found")
	ErrVersionAlreadyExists = errors.New("version already exists")
	ErrDatabaseErr          = errors.New("query error, please contact admin")
)

type IService interface {
	Create(ctx context.Context, req dto.CreatePromptVersionDTO) (*model.PromptVersion, error)
	Update(ctx context.Context, v *model.PromptVersion) error
	GetByID(ctx context.Context, id string) (*model.PromptVersion, error)
	GetByPromptID(ctx context.Context, promptID string) ([]*model.PromptVersion, error)
	GetLatestByPromptID(ctx context.Context, promptID string) (*model.PromptVersion, error)
	GetByPromptIDAndVersion(ctx context.Context, promptID, version string) (*model.PromptVersion, error)
	List(ctx context.Context, offset, limit int) ([]*model.PromptVersion, int64, error)
	DeleteByID(ctx context.Context, id string) error
}

type Service struct {
	repo       *version.Repo
	promptRepo *prompt.Repo
	logger     *zap.Logger
}

func CreateVersionService(repo *version.Repo, promptRepo *prompt.Repo, logger *zap.Logger) *Service {
	return &Service{
		repo:       repo,
		promptRepo: promptRepo,
		logger:     logger,
	}
}

func (s *Service) Create(ctx context.Context, req dto.CreatePromptVersionDTO) (*model.PromptVersion, error) {
	v := &model.PromptVersion{
		ID:        uuid.New().String(),
		PromptID:  req.PromptID,
		Version:   req.Version,
		Content:   req.Content,
		Variables: req.Variables,
		ChangeLog: req.ChangeLog,
		CreatedBy: req.CreatedBy,
		Username:  req.Username,
		IsPublish: req.IsPublish,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}

	// 如果发布版本，同步更新prompt原数据
	if req.IsPublish {
		if err := s.updatePromptMeta(ctx, v.PromptID, v.ID, v.IsPublish); err != nil {
			s.logger.Error("failed to update prompt meta: " + err.Error())
			// 不返回错误，因为version已创建成功
		}
	}

	return v, nil
}

// updatePromptMeta 更新prompt原数据的latestVersion和isPublish
func (s *Service) updatePromptMeta(ctx context.Context, promptID string, version string, isPublish bool) error {
	p, err := s.promptRepo.GetByID(ctx, promptID)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	if p == nil {
		return ErrPromptNotFound
	}

	p.LatestVersion = version
	p.IsPublish = isPublish
	return s.promptRepo.Update(ctx, p)
}

func (s *Service) Update(ctx context.Context, v *model.PromptVersion) error {
	old, err := s.repo.GetByID(ctx, v.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	if old == nil {
		return ErrVersionNotFound
	}

	v.BaseModel.CreatedAt = old.CreatedAt

	return s.repo.Update(ctx, v)
}

func (s *Service) GetByID(ctx context.Context, id string) (*model.PromptVersion, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return v, nil
}

func (s *Service) GetByPromptID(ctx context.Context, promptID string) ([]*model.PromptVersion, error) {
	v, err := s.repo.GetByPromptID(ctx, promptID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return v, nil
}

func (s *Service) GetLatestByPromptID(ctx context.Context, promptID string) (*model.PromptVersion, error) {
	v, err := s.repo.GetLatestByPromptID(ctx, promptID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}
	return v, nil
}

func (s *Service) GetByPromptIDAndVersion(ctx context.Context, promptID, version string) (*model.PromptVersion, error) {
	versions, err := s.repo.GetByPromptID(ctx, promptID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDatabaseErr
	}

	// 查找匹配的版本
	for _, v := range versions {
		if v.Version == version {
			return v, nil
		}
	}
	return nil, nil
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]*model.PromptVersion, int64, error) {
	v, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, 0, ErrDatabaseErr
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, 0, ErrDatabaseErr
	}
	return v, count, nil
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	old, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	if old == nil {
		return ErrVersionNotFound
	}

	if err := s.repo.DeleteByID(ctx, id); err != nil {
		s.logger.Error(err.Error())
		return ErrDatabaseErr
	}
	return nil
}
