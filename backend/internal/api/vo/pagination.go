package vo

// PageData 统一分页响应结构
type PageData[T any] struct {
	List  []*T  `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// NewPageData 创建分页数据
func NewPageData[T any](list []*T, total int64, page, limit int) *PageData[T] {
	return &PageData[T]{
		List:  list,
		Total: total,
		Page:  page,
		Limit: limit,
	}
}
