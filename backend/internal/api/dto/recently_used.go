package dto

type RecordUsageDTO struct {
	PromptID string `json:"promptId" binding:"required"`
}

type RemoveRecentlyUsedDTO struct {
	PromptID string `json:"promptId" binding:"required"`
}

type ListRecentlyUsedDTO struct {
	Limit int `json:"limit"`
}

type CleanRecentlyUsedDTO struct {
	KeepCount int `json:"keepCount"`
}
