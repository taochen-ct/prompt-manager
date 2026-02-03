package vo

import (
	"backend/internal/model"
	"backend/pkg/common"
)

type PromptVO struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	LatestVersion string `json:"latestVersion"`
	IsPublish     bool   `json:"isPublish"`
	CreateBy      string `json:"createBy"`
	Username      string `json:"username"`
	Category      string `json:"category"`
	CreateAt      string `json:"createAt"`
	UpdateAt      string `json:"updateAt"`
}

func FromPrompt(p *model.Prompt) *PromptVO {
	return &PromptVO{
		ID:            p.ID,
		Name:          p.Name,
		Path:          p.Path,
		LatestVersion: p.LatestVersion,
		IsPublish:     p.IsPublish,
		CreateBy:      p.CreatedBy,
		Username:      p.Username,
		Category:      p.Category,
		CreateAt:      common.FormatTime(p.CreatedAt),
		UpdateAt:      common.FormatTime(p.UpdatedAt),
	}
}

func FromPrompts(prompts []*model.Prompt) []*PromptVO {
	res := make([]*PromptVO, 0, len(prompts))
	for _, p := range prompts {
		res = append(res, FromPrompt(p))
	}
	return res
}
