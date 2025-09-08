package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type SeatRepository interface {
	GetAvailableSeats(scheduleID int) ([]models.AvailableSeat, error)
}

type seatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) SeatRepository {
	return &seatRepository{db}
}

func (r *seatRepository) GetAvailableSeats(scheduleID int) ([]models.AvailableSeat, error) {
	query := `
		SELECT s.id, s.cols, s.rows
		FROM seats s
		WHERE NOT EXISTS (
		    SELECT 1
		    FROM order_seats os
		    JOIN order_tickets ot ON os.order_id = ot.id
		    WHERE os.seat_id = s.id
		      AND ot.schedule_id = $1
		      AND ot.is_paid = true
		)
		ORDER BY s.cols, s.rows;
	`
	rows, err := r.db.Query(context.Background(), query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.AvailableSeat
	for rows.Next() {
		var seat models.AvailableSeat
		if err := rows.Scan(&seat.ID, &seat.Cols, &seat.Rows); err != nil {
			return nil, err
		}
		seat.Status = "available" // karena hanya kursi available yang diambil
		seats = append(seats, seat)
	}
	return seats, nil
}
