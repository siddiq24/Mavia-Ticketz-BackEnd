package repositories

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type MovieRepository interface {
	GetByFilter(name string, genre []int, page int) ([]models.Movie, int, error)
	CreateMovieWithSchedules(ctx context.Context, movieReq models.MovieRequest) (int, error)
	GetMovieByID(ctx context.Context, id int) (models.Movie, error)
	GetGenres(ctx context.Context) ([]models.Genre, error)
	GetUpcomingMovie(ctx context.Context) ([]models.Movie, error)
	GetPopularMovie(ctx context.Context) ([]models.Movie, error)
	DeleteMovieById(ctx context.Context, id int64) error
}

type movieRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) MovieRepository {
	return &movieRepository{db: db, rdb: rdb}
}

func (r *movieRepository) GetByFilter(name string, genre []int, page int) ([]models.Movie, int, error) {
	limit := 12
	offset := (page - 1) * limit

	// Base query dengan parameter yang dinamis
	var query, queryTotal string
	var params []interface{}
	var totalParams []interface{}

	// Parameter 1: name
	params = append(params, name)
	totalParams = append(totalParams, name)

	// Build query berdasarkan ada/tidaknya genre filter
	if len(genre) > 0 {
		query = `
			SELECT m.id, m.title, m.poster_img, m.rating,
				   COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres, m.release_date
			FROM movies m
			LEFT JOIN genre_movie mg ON m.id = mg.movie_id
			LEFT JOIN genres g ON mg.genre_id = g.id
			WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
			  AND m.id IN (
				SELECT movie_id
				FROM genre_movie
				WHERE genre_id = ANY($2)
				GROUP BY movie_id
				HAVING COUNT(DISTINCT genre_id) = $3
			)
				AND is_deleted = false
			GROUP BY m.id
			ORDER BY m.release_date DESC
			LIMIT $4 OFFSET $5
		`

		queryTotal = `
			SELECT COUNT(DISTINCT m.id)
			FROM movies m
			WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
			  AND m.id IN (
				SELECT movie_id
				FROM genre_movie
				WHERE genre_id = ANY($2)
				AND is_deleted = false
				GROUP BY movie_id
				HAVING COUNT(DISTINCT genre_id) = $3
			)
		`

		// Tambahkan parameter untuk genre
		params = append(params, pq.Array(genre), len(genre), limit, offset)
		totalParams = append(totalParams, pq.Array(genre), len(genre))

	} else {
		query = `
			SELECT m.id, m.title, m.poster_img, m.rating,
				   COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres, m.release_date
			FROM movies m
			LEFT JOIN genre_movie mg ON m.id = mg.movie_id
			LEFT JOIN genres g ON mg.genre_id = g.id
			WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
				AND is_deleted = false
			GROUP BY m.id
			ORDER BY m.release_date DESC
			LIMIT $2 OFFSET $3
		`

		queryTotal = `
			SELECT COUNT(DISTINCT m.id)
			FROM movies m
			WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
				AND is_deleted = false
		`

		// Tambahkan parameter untuk limit dan offset
		params = append(params, limit, offset)
	}

	// ===== Hitung total count =====
	ctx := context.Background()
	var totalCount int
	err := r.db.QueryRow(ctx, queryTotal, totalParams...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(limit)))

	// ===== Query data =====
	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie

		err := rows.Scan(&m.ID, &m.Title, &m.PosterImg, &m.Rating, &m.Genres, &m.ReleaseDate)
		if err != nil {
			return nil, 0, err
		}

		movies = append(movies, m)
	}

	return movies, totalPage, nil
}

func (r *movieRepository) CreateMovieWithSchedules(ctx context.Context, movieReq models.MovieRequest) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Handle Director - Get or Create
	directorID, err := r.getOrCreateDirector(ctx, tx, movieReq.DirectorName)
	if err != nil {
		return 0, fmt.Errorf("failed to handle director: %w", err)
	}

	// 2. Insert movie
	var movieID int
	query := `INSERT INTO movies (title, description, release_date, duration, poster_img, director_id, backdrop_img, rating) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err = tx.QueryRow(ctx, query,
		movieReq.Title,
		movieReq.Description,
		movieReq.ReleaseDate,
		movieReq.Duration,
		movieReq.PosterImg,
		directorID, // Use the resolved director ID
		movieReq.BackdropImg,
		movieReq.Rating,
	).Scan(&movieID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert movie: %w", err)
	}

	// 3. Insert genres
	if len(movieReq.GenreIDs) > 0 {
		genreQuery := "INSERT INTO genre_movie (movie_id, genre_id) VALUES ($1, $2)"
		batch := &pgx.Batch{}
		for _, genreID := range movieReq.GenreIDs {
			batch.Queue(genreQuery, movieID, genreID)
		}

		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return 0, fmt.Errorf("failed to insert genres: %w", err)
		}
	}

	// 4. Insert cast
	if len(movieReq.CasterIDs) > 0 {
		castQuery := "INSERT INTO caster_movie (movie_id, caster_id) VALUES ($1, $2)"
		batch := &pgx.Batch{}
		for _, casterID := range movieReq.CasterIDs {
			batch.Queue(castQuery, movieID, casterID)
		}

		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return 0, fmt.Errorf("failed to insert cast: %w", err)
		}
	}

	// 5. Insert schedules
	if len(movieReq.Schedules) > 0 {
		scheduleQuery := "INSERT INTO schedules (movie_id, cinema_id, time_id, date, city_id) VALUES ($1, $2, $3, $4, $5)"
		batch := &pgx.Batch{}

		for _, schedule := range movieReq.Schedules {
			// Parse schedule date
			scheduleDate, err := time.Parse("2006-01-02", schedule.Date)
			if err != nil {
				return 0, fmt.Errorf("invalid schedule date format: %w", err)
			}

			batch.Queue(scheduleQuery,
				movieID,
				schedule.CinemaID,
				schedule.TimeID,
				scheduleDate,
				schedule.CityID,
			)
		}

		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return 0, fmt.Errorf("failed to insert schedules: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return movieID, nil
}

func (r *movieRepository) getOrCreateDirector(ctx context.Context, tx pgx.Tx, directorName string) (int, error) {
	var directorID int

	// First, try to find existing director
	err := tx.QueryRow(ctx,
		"SELECT id FROM directors WHERE name = $1",
		directorName,
	).Scan(&directorID)

	if err == nil {
		// Director found, return the ID
		return directorID, nil
	}

	if err != pgx.ErrNoRows {
		// Some other error occurred
		return 0, fmt.Errorf("failed to query director: %w", err)
	}

	// Director not found, create new one
	err = tx.QueryRow(ctx,
		"INSERT INTO directors (name) VALUES ($1) RETURNING id",
		directorName,
	).Scan(&directorID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert director: %w", err)
	}

	return directorID, nil
}

func (r *movieRepository) GetMovieByID(ctx context.Context, id int) (models.Movie, error) {
	var movie models.Movie
	query := `SELECT m.id, m.title, m.description, m.release_date, m.duration, m.poster_img, d.name, m.backdrop_img, m.rating, m.created_at,
			  COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres,
			  COALESCE(array_agg(DISTINCT c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS casters
	          FROM movies m
			  LEFT JOIN genre_movie mg ON m.id = mg.movie_id
			  LEFT JOIN genres g ON mg.genre_id = g.id
			  LEFT JOIN caster_movie cm ON m.id = cm.movie_id
			  LEFT JOIN casters c ON cm.caster_id = c.id
			  LEFT JOIN directors d ON m.director_id = d.id
			  WHERE m.id = $1
			  GROUP BY m.id, d.name
			  ORDER BY m.id DESC`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Duration,
		&movie.PosterImg,
		&movie.Director,
		&movie.BackdropImg,
		&movie.Rating,
		&movie.CreatedAt,
		&movie.Genres,
		&movie.Cast,
	)
	if err != nil {
		return models.Movie{}, fmt.Errorf("failed to get movie: %w", err)
	}
	return movie, nil
}

func (r *movieRepository) GetGenres(ctx context.Context) ([]models.Genre, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM genres`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *movieRepository) GetUpcomingMovie(ctx context.Context) ([]models.Movie, error) {
	query := `
			SELECT m.id, m.title, m.poster_img, m.rating,
				   COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
			FROM movies m
			LEFT JOIN genre_movie mg ON m.id = mg.movie_id
			LEFT JOIN genres g ON mg.genre_id = g.id
			WHERE m.release_date > now()
			GROUP BY m.id
			ORDER BY m.release_date DESC
			LIMIT 12 OFFSET 0
		`

	// ===== Query data =====
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie

		err := rows.Scan(&m.ID, &m.Title, &m.PosterImg, &m.Rating, &m.Genres)
		if err != nil {
			return nil, err
		}

		movies = append(movies, m)
	}

	return movies, nil
}

func (r *movieRepository) GetPopularMovie(ctx context.Context) ([]models.Movie, error) {
	query := `
			SELECT m.id, m.title, m.poster_img, m.rating,
				   COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
			FROM movies m
			LEFT JOIN genre_movie mg ON m.id = mg.movie_id
			LEFT JOIN genres g ON mg.genre_id = g.id
            GROUP BY m.id
			ORDER BY m.rating DESC
			LIMIT 12
		`

	// ===== Query data =====
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie

		err := rows.Scan(&m.ID, &m.Title, &m.PosterImg, &m.Rating, &m.Genres)
		if err != nil {
			return nil, err
		}

		movies = append(movies, m)
	}

	return movies, nil
}

func (r *movieRepository) DeleteMovieById(ctx context.Context, id int64) error {
	query := `
	        UPDATE movies 
        SET is_deleted = TRUE
        WHERE id = $1 AND is_deleted = FALSE
        RETURNING id;`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}
