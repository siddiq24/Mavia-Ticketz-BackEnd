package models

import "time"

type Profile struct {
	UserID    int       `json:"user_id"`
	Avatar    string    `json:"avatar"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Birthdate time.Time `json:"birthdate"`
}

type UpdateProfileRequest struct {
	Avatar    *string `json:"avatar"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	Birthdate *string `json:"birthdate"`
}
