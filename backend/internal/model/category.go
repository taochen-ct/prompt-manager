package model

// Category 对应 prompt_categories 表（提示词分类）
type Category struct {
	ID    string `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Icon  string `json:"icon" db:"icon"`
	Count int    `json:"count" db:"count"`
	URL   string `json:"url" db:"url"`
	BaseModel
}

func (Category) TableName() string {
	return "prompt_categories"
}
