package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type DashboardRepository interface {
	GetTicketSalesByMovie(ctx context.Context, year int) ([]models.TicketSalesByMovie, error)
	GetTicketSalesByCategoryLocation(ctx context.Context, year int) ([]models.TicketSalesByCategoryLocation, error)
}

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{db: db}
}

// Sales per movie per week
func (r *dashboardRepository) GetTicketSalesByMovie(ctx context.Context, year int) ([]models.TicketSalesByMovie, error) {
	query := `
		SELECT m.title, EXTRACT(WEEK FROM t.sold_at)::int AS week,
		EXTRACT(MONTH FROM t.sold_at)::int AS month,
		EXTRACT(YEAR FROM t.sold_at)::int AS year,
		COUNT(*) AS total_sales
		FROM tickets t
		JOIN movies m ON t.movie_id = m.id
		WHERE EXTRACT(YEAR FROM t.sold_at) = $1
		GROUP BY m.title, week, month, year
		ORDER BY month, week;
	`

	rows, err := r.db.Query(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TicketSalesByMovie
	for rows.Next() {
		var t models.TicketSalesByMovie
		if err := rows.Scan(&t.MovieTitle, &t.Week, &t.Month, &t.Year, &t.TotalSales); err != nil {
			log.Println(err)
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

// Sales per category & location
func (r *dashboardRepository) GetTicketSalesByCategoryLocation(ctx context.Context, year int) ([]models.TicketSalesByCategoryLocation, error) {
	query := `
		SELECT c.name AS category, l.name AS location,
		EXTRACT(MONTH FROM t.sold_at)::int AS month,
		EXTRACT(YEAR FROM t.sold_at)::int AS year,
		COUNT(*) AS total_sales
		FROM tickets t
		JOIN movies m ON t.movie_id = m.id
		JOIN categories c ON m.category_id = c.id
		JOIN locations l ON t.location_id = l.id
		WHERE EXTRACT(YEAR FROM t.sold_at) = $1
		GROUP BY c.name, l.name, month, year
		ORDER BY month;
	`

	rows, err := r.db.Query(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TicketSalesByCategoryLocation
	for rows.Next() {
		var t models.TicketSalesByCategoryLocation
		if err := rows.Scan(&t.Category, &t.Location, &t.Month, &t.Year, &t.TotalSales); err != nil {
			log.Println(err)
			continue
		}
		result = append(result, t)
	}
	return result, nil
}
