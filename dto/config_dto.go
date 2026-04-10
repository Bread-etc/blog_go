package dto

// ===== Requests (请求入参) =====

type SaveConfigReq struct {
	Title       string `json:"title" binding:"required"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	GithubURL   string `json:"github_url"`
}

// ===== Response（响应出参）=====

type ConfigResp struct {
	Title       string `json:"title" binding:"required"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	GithubURL   string `json:"github_url"`
}
