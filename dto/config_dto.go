package dto

// ===== Requests (请求入参) =====

type SaveConfigReq struct {
	Title       string `json:"title" binding:"required,max=100"`
	Subtitle    string `json:"subtitle" binding:"max=255"`
	Description string `json:"description" binding:"max=1000"`
	Keywords    string `json:"keywords" binding:"max=255"`
	Author      string `json:"author" binding:"max=50"`
	Email       string `json:"email" binding:"omitempty,email,max=100"`
	GithubURL   string `json:"githubUrl" binding:"omitempty,url,max=255"`
}

// ===== Response（响应出参）=====

type ConfigResp struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	GithubURL   string `json:"githubUrl"`
}
