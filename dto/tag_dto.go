package dto

// ===== Requests (请求入参) =====

type CreateTagReq struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=100"`
}

type UpdateTagReq struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=100"`
}

// ===== Response (响应出参) =====
// 复用 common_dto 中的 TagBrief 作为返回值
