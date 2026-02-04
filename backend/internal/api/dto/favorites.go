package dto

type AddFavoriteDTO struct {
	PromptID string `json:"promptId" binding:"required"`
}

type RemoveFavoriteDTO struct {
	PromptID string `json:"promptId" binding:"required"`
}

type CheckFavoriteDTO struct {
	PromptID string `json:"promptId" binding:"required"`
}
