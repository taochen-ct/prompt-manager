package prompt

import (
	"backend/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type IRepo interface {
	Create(ctx context.Context, p *model.Prompt) error
	Update(ctx context.Context, p *model.Prompt) error
	GetByID(ctx context.Context, id string) (*model.Prompt, error)
	GetByPath(ctx context.Context, path string) (*model.Prompt, error)
	List(ctx context.Context, userID string, offset, limit int) ([]*model.Prompt, error)
	Count(ctx context.Context) (int64, error)
	CountByUser(ctx context.Context, userID string) (int64, error)
	DeleteByID(ctx context.Context, id string) error
}

type Repo struct {
	db *sqlx.DB
}

func CreatePromptRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, p *model.Prompt) (*model.Prompt, error) {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.IsPublish = false
	p.LatestVersion = ""
	query := `
		INSERT INTO prompt (
			id, name, path, latest_version, is_publish,
			created_by, username, category, created_at, updated_at
		) VALUES (
			:id, :name, :path, :latest_version, :is_publish,
			:created_by, :username, :category, :created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repo) Update(ctx context.Context, p *model.Prompt) error {
	const query = `
		UPDATE prompt
		SET
			name = :name,
			latest_version = :latest_version,
			is_publish = :is_publish,
			category = :category,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, p)
	return err
}

func (r *Repo) GetByID(ctx context.Context, id string) (*model.Prompt, error) {
	const query = `
		SELECT id, name, path, latest_version, is_publish, created_at, updated_at, created_by, username, category
		FROM prompt
		WHERE id = ?
	`
	var p model.Prompt
	err := r.db.GetContext(ctx, &p, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &p, err
}

func (r *Repo) GetByPath(ctx context.Context, path string) (*model.Prompt, error) {
	const query = `
		SELECT id, name, path, latest_version, is_publish, created_at, updated_at, created_by, username, category
		FROM prompt
		WHERE path = ?
	`
	var p model.Prompt
	err := r.db.GetContext(ctx, &p, query, path)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &p, err
}

func (r *Repo) List(ctx context.Context, userID string, offset, limit int) ([]*model.Prompt, error) {
	query := `
		SELECT id, name, path, latest_version, is_publish, created_at, updated_at, created_by, username, category
		FROM prompt
		WHERE created_by = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	var list []*model.Prompt
	err := r.db.SelectContext(ctx, &list, query, userID, limit, offset)
	return list, err
}

func (r *Repo) Count(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(1) FROM prompt`
	var count int64
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *Repo) CountByUser(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(1) FROM prompt WHERE created_by = ?`
	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	const query = `
		DELETE FROM prompt
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
