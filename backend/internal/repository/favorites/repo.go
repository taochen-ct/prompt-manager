package favorites

import (
	"backend/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type IRepo interface {
	Create(ctx context.Context, f *model.Favorite) error
	DeleteByUserAndPrompt(ctx context.Context, userID int64, promptID string) error
	GetByUserAndPrompt(ctx context.Context, userID int64, promptID string) (*model.Favorite, error)
	ListByUser(ctx context.Context, userID int64, offset, limit int) ([]*model.Favorite, error)
	ListByUserWithPrompt(ctx context.Context, userID int64, offset, limit int) ([]*model.FavoriteWithPrompt, error)
	CountByUser(ctx context.Context, userID int64) (int64, error)
	Exists(ctx context.Context, userID int64, promptID string) (bool, error)
}

type Repo struct {
	db *sqlx.DB
}

func CreateFavoriteRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, f *model.Favorite) error {
	now := time.Now()
	f.CreatedAt = now
	query := `
		INSERT INTO favorites (id, user_id, prompt_id, created_at)
		VALUES (:id, :user_id, :prompt_id, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, f)
	return err
}

func (r *Repo) DeleteByUserAndPrompt(ctx context.Context, userID int64, promptID string) error {
	query := `
		DELETE FROM favorites
		WHERE user_id = ? AND prompt_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, userID, promptID)
	return err
}

func (r *Repo) GetByUserAndPrompt(ctx context.Context, userID int64, promptID string) (*model.Favorite, error) {
	query := `
		SELECT id, user_id, prompt_id, created_at
		FROM favorites
		WHERE user_id = ? AND prompt_id = ?
	`
	var f model.Favorite
	err := r.db.GetContext(ctx, &f, query, userID, promptID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &f, err
}

func (r *Repo) ListByUser(ctx context.Context, userID int64, offset, limit int) ([]*model.Favorite, error) {
	query := `
		SELECT id, user_id, prompt_id, created_at
		FROM favorites
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	var list []*model.Favorite
	err := r.db.SelectContext(ctx, &list, query, userID, limit, offset)
	return list, err
}

func (r *Repo) ListByUserWithPrompt(ctx context.Context, userID int64, offset, limit int) ([]*model.FavoriteWithPrompt, error) {
	query := `
		SELECT
			f.id, f.user_id, f.prompt_id, f.created_at,
			p.name AS prompt_name, p.path AS prompt_path,
			p.latest_version AS prompt_version, p.category AS prompt_category
		FROM favorites f
		INNER JOIN prompt p ON f.prompt_id = p.id
		WHERE f.user_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`
	var list []*model.FavoriteWithPrompt
	err := r.db.SelectContext(ctx, &list, query, userID, limit, offset)
	return list, err
}

func (r *Repo) CountByUser(ctx context.Context, userID int64) (int64, error) {
	query := `
		SELECT COUNT(1) FROM favorites
		WHERE user_id = ?
	`
	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

func (r *Repo) Exists(ctx context.Context, userID int64, promptID string) (bool, error) {
	query := `
		SELECT COUNT(1) FROM favorites
		WHERE user_id = ? AND prompt_id = ?
	`
	var count int
	err := r.db.GetContext(ctx, &count, query, userID, promptID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
