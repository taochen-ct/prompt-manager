package vo

import (
	"backend/internal/model"
	"backend/pkg/common"
)

type FavoriteVO struct {
	ID        string `json:"id"`
	UserID    int64  `json:"userId"`
	PromptID  string `json:"promptId"`
	CreatedAt string `json:"createdAt"`
}

func FromFavorite(f *model.Favorite) *FavoriteVO {
	return &FavoriteVO{
		ID:        f.ID,
		UserID:    f.UserID,
		PromptID:  f.PromptID,
		CreatedAt: common.FormatTime(f.CreatedAt),
	}
}

func FromFavorites(list []*model.Favorite) []*FavoriteVO {
	res := make([]*FavoriteVO, 0, len(list))
	for _, f := range list {
		res = append(res, FromFavorite(f))
	}
	return res
}

type FavoriteWithPromptVO struct {
	FavoriteVO
	PromptName     string `json:"promptName"`
	PromptPath     string `json:"promptPath"`
	PromptVersion  string `json:"latestVersion"`
	PromptCategory string `json:"category"`
}

func FromFavoriteWithPrompt(f *model.FavoriteWithPrompt) *FavoriteWithPromptVO {
	return &FavoriteWithPromptVO{
		FavoriteVO:     *FromFavorite(&f.Favorite),
		PromptName:     f.PromptName,
		PromptPath:     f.PromptPath,
		PromptVersion:  f.PromptVersion,
		PromptCategory: f.PromptCategory,
	}
}

func FromFavoritesWithPrompt(list []*model.FavoriteWithPrompt) []*FavoriteWithPromptVO {
	res := make([]*FavoriteWithPromptVO, 0, len(list))
	for _, f := range list {
		res = append(res, FromFavoriteWithPrompt(f))
	}
	return res
}
