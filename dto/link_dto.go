package dto

// ===== Requests (请求入参) =====

type CreateLinkReq struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

type UpdateLinkReq struct {
	Name        string `json:"name"`
	URL         string `json:"url" binding:"omitempty,url"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

// ===== Response (响应出参) =====

type LinkResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}
