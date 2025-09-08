package models

import "time"

type CreateOrderRequest struct {
	ScheduleID      int      `json:"schedule_id" binding:"required"`
	Seats           []string `json:"seats" binding:"required"`
	Fullname        string   `json:"fullname" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	Phone           string   `json:"phone" binding:"required"`
	PaymentMethodID int      `json:"payment_method_id" binding:"required"`
}

type Order struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	ScheduleID      int       `json:"schedule_id"`
	Total           int       `json:"total"`
	Fullname        string    `json:"fullname"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	PaymentMethodID int       `json:"payment_method_id"`
	IsPaid          bool      `json:"is_paid"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type OrderHistory struct {
	OrderID     int       `json:"order_id"`
	MovieTitle  string    `json:"movie_title"`
	CinemaName  string    `json:"cinema_name"`
	Location    string    `json:"location"`
	ShowTime    string    `json:"show_time"`
	TotalAmount int       `json:"total_amount"`
	IsPaid      bool      `json:"is_paid"`
	CreatedAt   time.Time `json:"created_at"`
	Seats       string    `json:"seats"`
}
