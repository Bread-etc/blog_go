package dto

// ===== Requests (请求入参) =====

type CreateCategoryReq struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateCategoryReq struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ===== Response (响应出参) =====
// category较为简单，故直接复用common_dto中的 CategoryBrief 作为返回值
