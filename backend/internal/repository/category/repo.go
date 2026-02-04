package category

import (
	"backend/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type IRepo interface {
	Create(ctx context.Context, c *model.Category) error
	Update(ctx context.Context, c *model.Category) error
	GetByID(ctx context.Context, id string) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
	DeleteByID(ctx context.Context, id string) error
}

type Repo struct {
	db *sqlx.DB
}

func CreateCategoryRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, c *model.Category) (*model.Category, error) {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Count = 0
	query := `
		INSERT INTO categories (
			id, title, icon, count, url, username, created_by, created_at, updated_at
		) VALUES (
			:id, :title, :icon, :count, :url, :username, :created_by, :created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repo) Update(ctx context.Context, c *model.Category) error {
	const query = `
		UPDATE categories
		SET
			title = :title,
			icon = :icon,
			count = :count,
			url = :url,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, c)
	return err
}

func (r *Repo) GetByID(ctx context.Context, id string) (*model.Category, error) {
	const query = `
		SELECT id, title, icon, count, url, created_at, updated_at
		FROM categories
		WHERE id = ?
	`
	var c model.Category
	err := r.db.GetContext(ctx, &c, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r *Repo) List(ctx context.Context) ([]*model.Category, error) {
	const query = `
		SELECT id, title, icon, count, url, created_at, updated_at
		FROM categories
		ORDER BY created_at ASC
	`
	var list []*model.Category
	err := r.db.SelectContext(ctx, &list, query)
	return list, err
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	const query = `
		DELETE FROM categories
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
