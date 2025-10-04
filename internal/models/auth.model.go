package models

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" example:"123456"`
	NewPassword string `json:"new_password" example:"newpassword"`
}
