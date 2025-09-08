package models

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Duration    int       `json:"duration"`
	PosterImg   string    `json:"poster_img"`
	DirectorsID int       `json:"directors_id"`
	BackdropImg string    `json:"backdrop_img"`
	Rating      float32   `json:"rating"`
	IsUpcoming  bool      `json:"is_upcoming"`
	CreatedAt   time.Time `json:"created_at"`
	Genres      []string  `json:"genres"`
	Genre       string    `json:"genre"`
}

type UpdateMovieRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Duration    *int    `json:"duration"`
	ReleaseDate *string `json:"release_date"`
	Genre       *string `json:"genre"`
}
