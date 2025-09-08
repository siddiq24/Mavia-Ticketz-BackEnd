package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type OrderRepository interface {
	CreateOrder(userID int, req models.CreateOrderRequest) (*models.Order, error)
	GetOrderHistory(userID int) ([]models.OrderHistory, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(userID int, req models.CreateOrderRequest) (*models.Order, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	pricePerSeat := 10
	total := pricePerSeat * len(req.Seats)

	// insert order_tickets
	var order models.Order
	query := `
		INSERT INTO order_tickets (user_id, schedule_id, total_amount, fullname, email, phone, payment_method_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, schedule_id, total_amount, fullname, email, phone, payment_method_id, is_paid, is_active, create_at
	`
	err = tx.QueryRow(context.Background(), query,
		userID, req.ScheduleID, total,
		req.Fullname, req.Email, req.Phone, req.PaymentMethodID).
		Scan(&order.ID, &order.UserID, &order.ScheduleID, &order.Total,
			&order.Fullname, &order.Email, &order.Phone,
			&order.PaymentMethodID, &order.IsPaid, &order.IsActive, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	// cek setiap seat & insert order_seats
	for _, seat := range req.Seats {
		col := string(seat[0])
		row := seat[1:]

		var seatID int
		err := tx.QueryRow(context.Background(), `
			SELECT s.id
			FROM seats s
			WHERE s.cols=$1 AND s.rows=$2
			  AND NOT EXISTS (
				  SELECT 1 
				  FROM order_seats os
				  JOIN order_tickets ot ON os.order_id = ot.id
				  WHERE os.seat_id = s.id 
				    AND ot.schedule_id = $3
				    AND ot.is_paid = true
			  )
		`, col, row, req.ScheduleID).Scan(&seatID)
		if err != nil {
			return nil, fmt.Errorf("seat %s not available", seat)
		}

		// insert ke order_seats
		_, err = tx.Exec(context.Background(),
			`INSERT INTO order_seats (order_id, seat_id) VALUES ($1, $2)`,
			order.ID, seatID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetOrderHistory(userID int) ([]models.OrderHistory, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT 
			ot.id, m.title, c.name, l.name, t.time, 
			ot.total_amount, ot.is_paid, ot.create_at,
			STRING_AGG(CONCAT(s.cols, s.rows), ', ') AS seats
		FROM order_tickets ot
		JOIN schedules sc ON ot.schedule_id = sc.id
		JOIN movies m ON sc.movie_id = m.id
		JOIN cinemas c ON sc.cinema_id = c.id
		JOIN locations l ON c.location_id = l.id
		JOIN times t ON sc.time_id = t.id
		JOIN order_seats os ON ot.id = os.order_id
		JOIN seats s ON os.seat_id = s.id
		WHERE ot.user_id = $1
		GROUP BY ot.id, m.title, c.name, l.name, t.time, ot.total_amount, ot.is_paid, ot.create_at
		ORDER BY ot.create_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderHistory
	for rows.Next() {
		var h models.OrderHistory
		err := rows.Scan(
			&h.OrderID, &h.MovieTitle, &h.CinemaName, &h.Location, &h.ShowTime,
			&h.TotalAmount, &h.IsPaid, &h.CreatedAt, &h.Seats,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}
