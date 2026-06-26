package dto

// ===== Requests (请求入参) =====

type LoginReq struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ===== Response (响应出参) =====

type LoginResp struct {
	Token string    `json:"token"`
	User  UserBrief `json:"user"`
}

type PublicKeyResp struct {
	PublicKey string `json:"publicKey"`
}

type UserBrief struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
