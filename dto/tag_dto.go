package dto

// ===== Requests (请求入参) =====

type CreateTagReq struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateTagReq struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// ===== Response (响应出参) =====
// tag较为简单，故直接复用common_dto中的 TagBrief 作为返回值
