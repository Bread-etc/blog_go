package dto

// AuthorBrief 作者简要基本信息
type AuthorBrief struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// CategoryBrief 分类基本信息（删除无关的时间戳）
type CategoryBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagBrief 标签基本信息
type TagBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// 统一的分页响应外壳
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"size"`
}
