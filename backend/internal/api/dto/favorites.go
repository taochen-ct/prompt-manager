package dto

type AddFavoriteDTO struct {
	UserID   int64  `json:"userId" binding:"required"`
	PromptID string `json:"promptId" binding:"required"`
}

type RemoveFavoriteDTO struct {
	UserID   int64  `json:"userId" binding:"required"`
	PromptID string `json:"promptId" binding:"required"`
}

type CheckFavoriteDTO struct {
	UserID   int64  `json:"userId" binding:"required"`
	PromptID string `json:"promptId" binding:"required"`
}

type ListFavoritesDTO struct {
	UserID int64 `json:"userId" binding:"required"`
}
