package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type MovieRepository interface {
	GetUpcoming() ([]models.Movie, error)
	GetPopular(limit int) ([]models.Movie, error)
	GetByFilter(name, genre string, page int) ([]models.Movie, error)
	GetMovieByID(id int) (*models.Movie, error)
	GetAllMovies() ([]models.Movie, error)
	UpdateMovie(id int, req models.UpdateMovieRequest) error
	DeleteMovie(id int) error
}

type movieRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) MovieRepository {
	return &movieRepository{db: db, rdb: rdb}
}

func (r *movieRepository) GetUpcoming() ([]models.Movie, error) {
	ctx := context.Background()
	redisKey := "movies:upcoming"

	// cek redis
	cmd := r.rdb.Get(ctx, redisKey)
	if r.rdb == nil {
		log.Println("Redis client belum diinisialisasi!")
	}
	if cmd.Err() == nil {
		// cache hit
		var cachedMovies []models.Movie
		cmdByte, err := cmd.Bytes()
		if err == nil {
			if err := json.Unmarshal(cmdByte, &cachedMovies); err == nil {
				if len(cachedMovies) > 0 {
					log.Println("Cache hit: upcoming movies")
					return cachedMovies, nil
				}
			} else {
				log.Println("Unmarshal error:", err)
			}
		} else {
			log.Println("Redis GET bytes error:", err)
		}
	} else if cmd.Err() != redis.Nil {
		log.Println("Redis GET error:", cmd.Err())
	}

	// kalau cache miss → ambil dari DB
	query := `
        SELECT m.id, m.title, m.description, m.release_date, m.duration,
               m.poster_img, m.backdrop_img, m.rating, m.is_upcoming,
               m.created_at, m.directors_id,
               COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
        FROM movies m
        LEFT JOIN genre_movie mg ON m.id = mg.movie_id
        LEFT JOIN genres g ON mg.genre_id = g.id
        WHERE m.is_upcoming = true
        GROUP BY m.id
        ORDER BY m.release_date ASC
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration,
			&m.PosterImg, &m.BackdropImg, &m.Rating, &m.IsUpcoming,
			&m.CreatedAt, &m.DirectorsID, &m.Genres); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	// simpan ke cache
	if len(movies) > 0 {
		bt, err := json.Marshal(movies)
		if err == nil {
			if err := r.rdb.Set(ctx, redisKey, string(bt), 24*time.Hour).Err(); err != nil {
				log.Println("Redis SET error:", err)
			} else {
				log.Println("Cache updated: upcoming movies")
			}
		} else {
			log.Println("Marshal error:", err)
		}
	}
	if r.rdb == nil {
		log.Println("Redis client belum diinisialisasi!")
	}

	return movies, nil
}

func (r *movieRepository) GetPopular(limit int) ([]models.Movie, error) {
	rows, err := r.db.Query(context.Background(),
		`
		SELECT m.id, m.title, m.description, m.release_date, m.duration,
		m.poster_img, m.backdrop_img, m.rating, m.is_upcoming,
		m.created_at, m.directors_id,
		COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
		FROM movies m
		LEFT JOIN genre_movie mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		GROUP BY m.id
		ORDER BY m.rating DESC
		LIMIT $1
		`, limit)
	if err != nil {
		log.Println("ERROR QUERY GetUpcoming:", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration,
			&m.PosterImg, &m.BackdropImg, &m.Rating, &m.IsUpcoming,
			&m.CreatedAt, &m.DirectorsID, &m.Genres)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *movieRepository) GetByFilter(name, genre string, page int) ([]models.Movie, error) {
	limit := 20
	offset := (page - 1) * limit
	query := `
		SELECT m.id, m.title, m.description, m.release_date, m.duration,
		       m.poster_img, m.backdrop_img, m.rating, m.is_upcoming,
		       m.created_at, m.directors_id,
		       COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
		FROM movies m
		LEFT JOIN genre_movie mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR g.name ILIKE '%' || $2 || '%')
		GROUP BY m.id
		ORDER BY m.release_date DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(context.Background(), query, name, genre, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration,
			&m.PosterImg, &m.BackdropImg, &m.Rating, &m.IsUpcoming,
			&m.CreatedAt, &m.DirectorsID, &m.Genres)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *movieRepository) GetMovieByID(id int) (*models.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.release_date, m.duration,
		       m.poster_img, m.backdrop_img, m.rating, m.is_upcoming,
		       m.created_at, m.directors_id,
		       COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
		FROM movies m
		LEFT JOIN genre_movie mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE m.id = $1
		GROUP BY m.id
	`
	row := r.db.QueryRow(context.Background(), query, id)

	var m models.Movie
	err := row.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration,
		&m.PosterImg, &m.BackdropImg, &m.Rating, &m.IsUpcoming,
		&m.CreatedAt, &m.DirectorsID, &m.Genres)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, err
	}

	return &m, nil
}

func (r *movieRepository) GetAllMovies() ([]models.Movie, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, title, description, duration, release_date, genre, created_at FROM movies ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Duration, &m.ReleaseDate, &m.Genre, &m.CreatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *movieRepository) UpdateMovie(id int, req models.UpdateMovieRequest) error {
	query := `
		UPDATE movies 
		SET title = COALESCE($1, title),
		    description = COALESCE($2, description),
		    duration = COALESCE($3, duration),
		    release_date = COALESCE($4, release_date),
		    genre = COALESCE($5, genre)
		WHERE id = $6
	`
	_, err := r.db.Exec(context.Background(), query,
		req.Title, req.Description, req.Duration, req.ReleaseDate, req.Genre, id)
	return err
}

func (r *movieRepository) DeleteMovie(id int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM movies WHERE id=$1", id) // nanti diganti
	return err
}
