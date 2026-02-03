package vo

import (
	"backend/internal/model"
	"backend/pkg/common"
)

type PromptVersionVO struct {
	ID        string `json:"id"`
	PromptID  string `json:"promptId"`
	Version   string `json:"version"`
	Content   string `json:"content"`
	Variables string `json:"variables"`
	IsPublish bool   `json:"isPublish"`
	ChangeLog string `json:"changeLog"`
	CreatedBy string `json:"createdBy"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func FromPromptVersion(v *model.PromptVersion) *PromptVersionVO {
	if v == nil {
		return nil
	}
	return &PromptVersionVO{
		ID:        v.ID,
		PromptID:  v.PromptID,
		Version:   v.Version,
		Content:   v.Content,
		Variables: v.Variables,
		IsPublish: v.IsPublish,
		ChangeLog: v.ChangeLog,
		CreatedBy: v.CreatedBy,
		Username:  v.Username,
		CreatedAt: common.FormatTime(v.CreatedAt),
		UpdatedAt: common.FormatTime(v.UpdatedAt),
	}
}

func FromPromptVersions(versions []*model.PromptVersion) []*PromptVersionVO {
	res := make([]*PromptVersionVO, 0, len(versions))
	for _, v := range versions {
		res = append(res, FromPromptVersion(v))
	}
	return res
}
