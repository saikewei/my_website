package auth

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}
