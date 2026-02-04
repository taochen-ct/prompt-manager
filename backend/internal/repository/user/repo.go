package user

import (
	"backend/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type IRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	SoftDelete(ctx context.Context, id int64) error
}

type Repo struct {
	db *sqlx.DB
}

func CreateRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, user *model.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsDeleted = false

	query := `
		INSERT INTO users (
			username, password, nickname, department,
			created_at, updated_at, is_deleted
		) VALUES (
			:username, :password, :nickname, :department,
			:created_at, :updated_at, :is_deleted
		)
	`

	res, err := r.db.NamedExecContext(ctx, query, user)
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return err
}

func (r *Repo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User

	query := `
		SELECT *
		FROM users
		WHERE id = ? AND is_deleted = 0
	`

	err := r.db.GetContext(ctx, &user, r.rebind(query), id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	query := `
		SELECT *
		FROM users
		WHERE username = ? AND is_deleted = 0
	`

	err := r.db.GetContext(ctx, &user, r.rebind(query), username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) List(ctx context.Context, offset, limit int) ([]*model.User, error) {
	users := make([]*model.User, 0)

	query := `
		SELECT id, username, nickname, department, created_at, updated_at, is_deleted
		FROM users
		WHERE is_deleted = 0
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(
		ctx,
		&users,
		r.rebind(query),
		limit,
		offset,
	)
	return users, err
}

func (r *Repo) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users SET
			nickname = :nickname,
			department = :department,
			updated_at = :updated_at
		WHERE id = :id AND is_deleted = 0
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *Repo) SoftDelete(ctx context.Context, id int64) error {
	query := `
		UPDATE users SET
			is_deleted = 1,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		r.rebind(query),
		time.Now(),
		id,
	)
	return err
}

func (r *Repo) rebind(query string) string {
	return r.db.Rebind(query)
}
