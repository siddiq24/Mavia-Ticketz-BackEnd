package models

import "time"

type Profile struct {
	UserID    int       `json:"user_id" example:"24"`
	Username  string    `json:"username" example:"Martin Paes"`
	Phone     string    `json:"phone" example:"0891237654"`
	Point     int       `json:"point" example:"86"`
	Avatar    string    `json:"avatar" example:"/uploads/avatars/ucup.png"`
	Address   string    `json:"address" example:"jl.situmorang"`
	Birthdate time.Time `json:"-" example:"2000-06-24"`
	Role      string    `json:"role" example:"user"`
	Token     string    `json:"token" example:"yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1LCJ1c2VybmFtZSI6IlNpdGkgQW1pbmFoIiwicm9sZSI6InVzZXIiLCJleHAiOjE3NTgxMjAzNDMsImlhdCI6MTc1ODAzMzk0M30.KBHwCpiUx4pVUbflSYYTIM3fIbaMOXL8e4RGCVhRmHk"`
	Password  string    `json:"-"`
	Email     string    `json:"-"`
}

type UpdateProfileRequest struct {
	Avatar    *string `json:"avatar"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	Birthdate *string `json:"birthdate"`
}
