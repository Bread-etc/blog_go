package dto

// ===== Requests (请求入参) =====

type CreateLinkReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	URL         string `json:"url" binding:"required,url,max=255"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort" binding:"min=0"`
}

type UpdateLinkReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	URL         string `json:"url" binding:"required,url,max=255"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort" binding:"min=0"`
}

// ===== Response (响应出参) =====

type LinkResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}
