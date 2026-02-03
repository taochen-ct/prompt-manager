package vo

import (
	"backend/internal/model"
	"backend/pkg/common"
)

type CategoryVO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Count     int    `json:"count"`
	URL       string `json:"url"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func FromCategory(c *model.Category) *CategoryVO {
	return &CategoryVO{
		ID:        c.ID,
		Title:     c.Title,
		Icon:      c.Icon,
		Count:     c.Count,
		URL:       c.URL,
		CreatedAt: common.FormatTime(c.CreatedAt),
		UpdatedAt: common.FormatTime(c.UpdatedAt),
	}
}

func FromCategories(list []*model.Category) []*CategoryVO {
	res := make([]*CategoryVO, 0, len(list))
	for _, c := range list {
		res = append(res, FromCategory(c))
	}
	return res
}
