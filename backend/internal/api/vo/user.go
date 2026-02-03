package vo

import (
	"backend/internal/model"
	"backend/pkg/common"
	"strconv"
)

type UserVO struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Department string `json:"department"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func FromUser(u *model.User) *UserVO {
	if u == nil {
		return nil
	}

	return &UserVO{
		ID:         strconv.FormatInt(u.ID, 10),
		Username:   u.Username,
		Nickname:   u.Nickname,
		Department: u.Department,
		CreatedAt:  common.FormatTime(u.CreatedAt),
		UpdatedAt:  common.FormatTime(u.UpdatedAt),
	}
}

func FromUsers(users []*model.User) []*UserVO {
	res := make([]*UserVO, 0, len(users))
	for _, u := range users {
		res = append(res, FromUser(u))
	}
	return res
}
