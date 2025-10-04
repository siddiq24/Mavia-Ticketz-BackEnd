package models

import "time"

type CreateOrderRequeste struct {
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
	CinemaImg   string    `json:"cinema_img"`
}

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MovieID   int       `json:"movie_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"` // pending, paid, canceled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	MovieID  int     `json:"movie_id"`
	Quantity int     `json:"quantity"`
	Total    float64 `json:"total"`
}

type OrderTicket struct {
	ID            int64  `json:"id"`
	TotalAmount   int    `json:"total_amount"`
	IsPaid        bool   `json:"is_paid"`
	PaymentMethod *int64 `json:"payment_method_id,omitempty"`
	ScheduleID    int64  `json:"schedule_id"`
	CinemaName    string `json:"cinema_name"`
	CinemaImage   string `json:"cinema_image"`
	Time          string `json:"time"`
	MovieTitle    string `json:"movie_title"`
	Date          string `json:"date"`
	SeatCol       string `json:"seat_col"`
	SeatRow       string `json:"seat_row"`
}

type CreateOrderRequest struct {
	PersonalInfo struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	} `json:"personalInfo"`
	Total      int      `json:"total"`
	Seats      []string `json:"seats"`
	ScheduleID int64    `json:"schedule_id"`
}
