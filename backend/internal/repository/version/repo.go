package version

import (
	"backend/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type IRepo interface {
	Create(ctx context.Context, v *model.PromptVersion) error
	Update(ctx context.Context, v *model.PromptVersion) error
	GetByID(ctx context.Context, id string) (*model.PromptVersion, error)
	GetByPromptID(ctx context.Context, promptID string) ([]*model.PromptVersion, error)
	GetLatestByPromptID(ctx context.Context, promptID string) (*model.PromptVersion, error)
	List(ctx context.Context, offset, limit int) ([]*model.PromptVersion, error)
	Count(ctx context.Context) (int64, error)
	DeleteByID(ctx context.Context, id string) error
}

type Repo struct {
	db *sqlx.DB
}

func CreateVersionRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, v *model.PromptVersion) error {
	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	query := `
		INSERT INTO prompt_version (
			id, prompt_id, version, content, variables,
			is_publish, change_log, created_by, username,
			created_at, updated_at
		) VALUES (
			:id, :prompt_id, :version, :content, :variables,
			:is_publish, :change_log, :created_by, :username,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, v)
	return err
}

func (r *Repo) Update(ctx context.Context, v *model.PromptVersion) error {
	v.UpdatedAt = time.Now()
	query := `
		UPDATE prompt_version SET
			version = :version,
			content = :content,
			variables = :variables,
			is_publish = :is_publish,
			change_log = :change_log,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, v)
	return err
}

func (r *Repo) GetByID(ctx context.Context, id string) (*model.PromptVersion, error) {
	const query = `
		SELECT id, prompt_id, version, content, variables,
			is_publish, change_log, created_by, username,
			created_at, updated_at
		FROM prompt_version
		WHERE id = ?
	`
	var v model.PromptVersion
	err := r.db.GetContext(ctx, &v, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &v, err
}

func (r *Repo) GetByPromptID(ctx context.Context, promptID string) ([]*model.PromptVersion, error) {
	const query = `
		SELECT id, prompt_id, version, content, variables,
			is_publish, change_log, created_by, username,
			created_at, updated_at
		FROM prompt_version
		WHERE prompt_id = ?
		ORDER BY created_at DESC
	`
	var list []*model.PromptVersion
	err := r.db.SelectContext(ctx, &list, query, promptID)
	return list, err
}

func (r *Repo) GetLatestByPromptID(ctx context.Context, promptID string) (*model.PromptVersion, error) {
	const query = `
		SELECT id, prompt_id, version, content, variables,
			is_publish, change_log, created_by, username,
			created_at, updated_at
		FROM prompt_version
		WHERE prompt_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var v model.PromptVersion
	err := r.db.GetContext(ctx, &v, query, promptID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &v, err
}

func (r *Repo) List(ctx context.Context, offset, limit int) ([]*model.PromptVersion, error) {
	const query = `
		SELECT id, prompt_id, version, content, variables,
			is_publish, change_log, created_by, username,
			created_at, updated_at
		FROM prompt_version
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	var list []*model.PromptVersion
	err := r.db.SelectContext(ctx, &list, query, limit, offset)
	return list, err
}

func (r *Repo) Count(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(1) FROM prompt_version`
	var count int64
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	const query = `
		DELETE FROM prompt_version
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repo) DeleteByPromptId(ctx context.Context, promptId string) error {
	const query = `
		DELETE FROM prompt_version
		WHERE prompt_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, promptId)
	return err
}
