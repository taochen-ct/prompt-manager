package dto

type CreatePromptDTO struct {
	Name      string `json:"name" binding:"required"`
	CreatedBy string `json:"createdBy" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Path      string `json:"path" binding:"required"`
	Category  string `json:"category"`
}

type UpdatePromptDTO struct {
	ID        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	IsPublish bool   `json:"isPublish"`
	Category  string `json:"category"`
}
