package recently_used

import (
	"backend/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type IRepo interface {
	CreateOrUpdate(ctx context.Context, r *model.RecentlyUsed) (*model.RecentlyUsed, error)
	DeleteByUserAndPrompt(ctx context.Context, userID int64, promptID string) error
	GetByUserAndPrompt(ctx context.Context, userID int64, promptID string) (*model.RecentlyUsed, error)
	ListByUser(ctx context.Context, userID int64, offset, limit int) ([]*model.RecentlyUsed, error)
	ListByUserWithPrompt(ctx context.Context, userID int64, offset, limit int) ([]*model.RecentlyUsedWithPrompt, error)
	CountByUser(ctx context.Context, userID int64) (int64, error)
	DeleteOldByUser(ctx context.Context, userID int64, keepCount int) error
}

type Repo struct {
	db *sqlx.DB
}

func CreateRecentlyUsedRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateOrUpdate(ctx context.Context, rec *model.RecentlyUsed) (*model.RecentlyUsed, error) {
	now := time.Now()
	rec.UsedAt = now

	// 先查询是否存在
	existing, err := r.GetByUserAndPrompt(ctx, rec.UserID, rec.PromptID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// 已存在，更新使用时间
		query := `
			UPDATE recently_used
			SET used_at = ?
			WHERE user_id = ? AND prompt_id = ?
		`
		_, err := r.db.ExecContext(ctx, query, now, rec.UserID, rec.PromptID)
		if err != nil {
			return nil, err
		}
		rec.ID = existing.ID
		return rec, nil
	}

	// 不存在，插入新记录
	query := `
		INSERT INTO recently_used (id, user_id, prompt_id, used_at)
		VALUES (:id, :user_id, :prompt_id, :used_at)
	`
	_, err = r.db.NamedExecContext(ctx, query, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *Repo) DeleteByUserAndPrompt(ctx context.Context, userID int64, promptID string) error {
	query := `
		DELETE FROM recently_used
		WHERE user_id = ? AND prompt_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, userID, promptID)
	return err
}

func (r *Repo) GetByUserAndPrompt(ctx context.Context, userID int64, promptID string) (*model.RecentlyUsed, error) {
	query := `
		SELECT id, user_id, prompt_id, used_at
		FROM recently_used
		WHERE user_id = ? AND prompt_id = ?
	`
	var rec model.RecentlyUsed
	err := r.db.GetContext(ctx, &rec, query, userID, promptID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &rec, err
}

func (r *Repo) ListByUser(ctx context.Context, userID int64, offset, limit int) ([]*model.RecentlyUsed, error) {
	query := `
		SELECT id, user_id, prompt_id, used_at
		FROM recently_used
		WHERE user_id = ?
		ORDER BY used_at DESC
		LIMIT ? OFFSET ?
	`
	var list []*model.RecentlyUsed
	err := r.db.SelectContext(ctx, &list, query, userID, limit, offset)
	return list, err
}

func (r *Repo) ListByUserWithPrompt(ctx context.Context, userID int64, offset, limit int) ([]*model.RecentlyUsedWithPrompt, error) {
	query := `
		SELECT
			r.id, r.user_id, r.prompt_id, r.used_at,
			p.name AS prompt_name, p.path AS prompt_path,
			p.latest_version AS prompt_version, p.category AS prompt_category
		FROM recently_used r
		INNER JOIN prompt p ON r.prompt_id = p.id
		WHERE r.user_id = ?
		ORDER BY r.used_at DESC
		LIMIT ? OFFSET ?
	`
	var list []*model.RecentlyUsedWithPrompt
	err := r.db.SelectContext(ctx, &list, query, userID, limit, offset)
	return list, err
}

func (r *Repo) CountByUser(ctx context.Context, userID int64) (int64, error) {
	query := `
		SELECT COUNT(1) FROM recently_used
		WHERE user_id = ?
	`
	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

func (r *Repo) DeleteOldByUser(ctx context.Context, userID int64, keepCount int) error {
	query := `
		DELETE FROM recently_used
		WHERE user_id = ? AND id NOT IN (
			SELECT id FROM recently_used
			WHERE user_id = ?
			ORDER BY used_at DESC
			LIMIT ?
		)
	`
	_, err := r.db.ExecContext(ctx, query, userID, userID, keepCount)
	return err
}
