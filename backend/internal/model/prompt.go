package model

// Prompt 对应 prompt 表（提示词元信息）
type Prompt struct {
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Path          string `json:"path" db:"path"`
	LatestVersion string `json:"latest_version" db:"latest_version"`
	IsPublish     bool   `json:"is_publish" db:"is_publish"`
	CreatedBy     string `json:"created_by" db:"created_by"`
	Username      string `json:"username" db:"username"`
	Category      string `json:"category" db:"category"`
	BaseModel
}

func (Prompt) TableName() string {
	return "prompt"
}

type PromptVersion struct {
	ID        string `json:"id" db:"id"`
	PromptID  string `json:"prompt_id" db:"prompt_id"`
	Version   string `json:"version" db:"version"`
	Content   string `json:"content" db:"content"`
	Variables string `json:"variables" db:"variables"`
	IsPublish bool   `json:"is_publish" db:"is_publish"`
	ChangeLog string `json:"change_log" db:"change_log"`
	CreatedBy string `json:"created_by" db:"created_by"`
	Username  string `json:"username" db:"username"`
	BaseModel
}

func (PromptVersion) TableName() string {
	return "prompt_version"
}
