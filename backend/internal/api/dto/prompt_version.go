package dto

type CreatePromptVersionDTO struct {
	PromptID  string `json:"promptId" binding:"required"`
	Version   string `json:"version" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Variables string `json:"variables"`
	ChangeLog string `json:"changeLog"`
	CreatedBy string `json:"createdBy" binding:"required"`
	Username  string `json:"username" binding:"required"`
	IsPublish bool   `json:"isPublish"`
}

type UpdatePromptVersionDTO struct {
	ID        string `json:"id" binding:"required"`
	Version   string `json:"version" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Variables string `json:"variables"`
	ChangeLog string `json:"changeLog"`
	IsPublish bool   `json:"isPublish"`
}

type ListPromptVersionDTO struct {
	PromptID string `json:"promptId"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}
