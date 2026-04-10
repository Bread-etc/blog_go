package dto

// ===== Requests (请求入参) =====

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ===== Response (响应出参) =====

// LoginResp 登录成功响应
type LoginResp struct {
	Token string    `json:"token"`
	User  UserBrief `json:"user"`
}

// UserBrief 用户安全信息 (不包含 Password、Email 等敏感字段)
type UserBrief struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
