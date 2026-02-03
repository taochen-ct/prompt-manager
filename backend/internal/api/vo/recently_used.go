package vo

import (
	"backend/internal/model"
	"backend/pkg/common"
)

type RecentlyUsedVO struct {
	ID       string `json:"id"`
	UserID   int64  `json:"userId"`
	PromptID string `json:"promptId"`
	UsedAt   string `json:"usedAt"`
}

func FromRecentlyUsed(r *model.RecentlyUsed) *RecentlyUsedVO {
	return &RecentlyUsedVO{
		ID:       r.ID,
		UserID:   r.UserID,
		PromptID: r.PromptID,
		UsedAt:   common.FormatTime(r.UsedAt),
	}
}

func FromRecentlyUsedList(list []*model.RecentlyUsed) []*RecentlyUsedVO {
	res := make([]*RecentlyUsedVO, 0, len(list))
	for _, r := range list {
		res = append(res, FromRecentlyUsed(r))
	}
	return res
}

type RecentlyUsedWithPromptVO struct {
	RecentlyUsedVO
	PromptName     string `json:"promptName"`
	PromptPath     string `json:"promptPath"`
	PromptVersion  string `json:"latestVersion"`
	PromptCategory string `json:"category"`
}

func FromRecentlyUsedWithPrompt(r *model.RecentlyUsedWithPrompt) *RecentlyUsedWithPromptVO {
	return &RecentlyUsedWithPromptVO{
		RecentlyUsedVO: *FromRecentlyUsed(&r.RecentlyUsed),
		PromptName:     r.PromptName,
		PromptPath:     r.PromptPath,
		PromptVersion:  r.PromptVersion,
		PromptCategory: r.PromptCategory,
	}
}

func FromRecentlyUsedWithPromptList(list []*model.RecentlyUsedWithPrompt) []*RecentlyUsedWithPromptVO {
	res := make([]*RecentlyUsedWithPromptVO, 0, len(list))
	for _, r := range list {
		res = append(res, FromRecentlyUsedWithPrompt(r))
	}
	return res
}
