package model

import "time"

// RecentlyUsed 对应 recently_used 表（最近使用）
type RecentlyUsed struct {
	ID       string    `json:"id" db:"id"`
	UserID   int64     `json:"userId" db:"user_id"`
	PromptID string    `json:"promptId" db:"prompt_id"`
	UsedAt   time.Time `db:"used_at" json:"usedAt"`
}

// RecentlyUsedWithPrompt 最近使用详情（含prompt信息）
type RecentlyUsedWithPrompt struct {
	RecentlyUsed
	PromptName     string `db:"prompt_name" json:"promptName"`
	PromptPath     string `db:"prompt_path" json:"promptPath"`
	PromptVersion  string `db:"prompt_version" json:"latestVersion"`
	PromptCategory string `db:"prompt_category" json:"category"`
}

func (RecentlyUsed) TableName() string {
	return "recently_used"
}
