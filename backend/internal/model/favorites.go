package model

import "time"

// Favorite 对应 favorites 表（收藏夹）
type Favorite struct {
	ID        string    `json:"id" db:"id"`
	UserID    int64     `json:"userId" db:"user_id"`
	PromptID  string    `json:"promptId" db:"prompt_id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// FavoriteWithPrompt 收藏夹详情（含prompt信息）
type FavoriteWithPrompt struct {
	Favorite
	PromptName     string `db:"prompt_name" json:"promptName"`
	PromptPath     string `db:"prompt_path" json:"promptPath"`
	PromptVersion  string `db:"prompt_version" json:"latestVersion"`
	PromptCategory string `db:"prompt_category" json:"category"`
}

func (Favorite) TableName() string {
	return "favorites"
}
