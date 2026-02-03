package model

type User struct {
	ID         int64  `db:"id" json:"id"`
	Username   string `db:"username" json:"username"`
	Nickname   string `db:"nickname" json:"nickname"`
	Department string `db:"department" json:"department"`
	IsDeleted  bool   `db:"is_deleted" json:"isDeleted"`
	BaseModel
}
