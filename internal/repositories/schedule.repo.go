package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type ScheduleRepository interface {
	GetSchedules() ([]models.Schedule, error)
}

type scheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) GetSchedules() ([]models.Schedule, error) {
	query := `
	SELECT s.id, s.date,
	       m.id, m.title, m.duration, m.poster_img, m.rating,
	       c.id, c.name, c.image,
	       l.id, l.name,
	       t.id, t.time
	FROM schedules s
	JOIN movies m ON s.movie_id = m.id
	JOIN cinemas c ON s.cinema_id = c.id
	JOIN locations l ON c.location_id = l.id
	JOIN times t ON s.time_id = t.id
	ORDER BY s.date, t.time;
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		err := rows.Scan(
			&s.ID, &s.Date,
			&s.Movie.ID, &s.Movie.Title, &s.Movie.Duration, &s.Movie.PosterImg, &s.Movie.Rating,
			&s.Cinema.ID, &s.Cinema.Name, &s.Cinema.Image,
			&s.Location.ID, &s.Location.Name,
			&s.Time.ID, &s.Time.Time,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}
