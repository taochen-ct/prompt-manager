package dto

type CreateCategoryDTO struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Icon  string `json:"icon" binding:"required"`
	URL   string `json:"url" binding:"required"`
}

type UpdateCategoryDTO struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Icon  string `json:"icon" binding:"required"`
	Count int    `json:"count"`
	URL   string `json:"url" binding:"required"`
}
