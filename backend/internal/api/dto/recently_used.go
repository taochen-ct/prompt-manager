package dto

type RecordUsageDTO struct {
	UserID   int64  `json:"userId" binding:"required"`
	PromptID string `json:"promptId" binding:"required"`
}

type RemoveRecentlyUsedDTO struct {
	UserID   int64  `json:"userId" binding:"required"`
	PromptID string `json:"promptId" binding:"required"`
}

type ListRecentlyUsedDTO struct {
	UserID int64 `json:"userId" binding:"required"`
	Limit  int   `json:"limit"`
}

type CleanRecentlyUsedDTO struct {
	UserID    int64 `json:"userId" binding:"required"`
	KeepCount int   `json:"keepCount"`
}
