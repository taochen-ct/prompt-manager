package dto

// CreateUserDTO 创建用户
type CreateUserDTO struct {
	Username   string `json:"username" binding:"required,min=3,max=32"`
	Nickname   string `json:"nickname"`
	Department string `json:"department"`
}

// UpdateUserDTO 更新用户
type UpdateUserDTO struct {
	Nickname   string `json:"nickname"`
	Department string `json:"department"`
}

type DeleteUserDTO struct {
	Id         string `json:"id" binding:"required"`
	Nickname   string `json:"nickname"`
	Department string `json:"department"`
}
