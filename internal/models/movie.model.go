package models

import (
	"mime/multipart"
	"time"
)

// type Movie struct {
// 	ID          int       `json:"id"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	ReleaseDate time.Time `json:"release_date"`
// 	Duration    int       `json:"duration"`
// 	PosterImg   string    `json:"poster_img"`
// 	DirectorsID int       `json:"directors_id"`
// 	BackdropImg string    `json:"backdrop_img"`
// 	Rating      float32   `json:"rating"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	Genres      []string  `json:"genres"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	CastersID   []int     `json:"caster_id"`
// }

type UpdateMovieRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Duration    *int    `json:"duration"`
	ReleaseDate *string `json:"release_date"`
	Genres      *[]int  `json:"genres"`
}

type CreateMovieRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	ReleaseDate string  `json:"release_date"`
	Duration    int     `json:"duration" binding:"required"`
	PosterImg   string  `json:"poster_img"`
	BackdropImg string  `json:"backdrop_img"`
	Rating      float64 `json:"rating"`
	DirectorsID int     `json:"directors_id" binding:"required"`
	CasterIDs   []int   `json:"caster_ids"`
	GenreIDs    []int   `json:"genre_ids"`
}

type MovieUploadBody struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	ReleaseDate string                `form:"release_date"`
	Duration    int                   `form:"duration"`
	DirectorsId int                   `form:"directors_id"`
	PosterImg   *multipart.FileHeader `form:"poster"`
	BackdropImg *multipart.FileHeader `form:"backdrop"`
}

type Genre struct {
	ID   int    `json:"id,omitempty" db:"id"`
	Name string `json:"name" db:"name"`
}

type Caster struct {
	ID   int    `json:"id,omitempty" db:"id"`
	Name string `json:"name" db:"name"`
}

type Movie struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
	Duration    int       `json:"duration,omitempty" db:"duration"`
	PosterImg   string    `json:"poster_img" db:"poster_img"`
	DirectorID  int       `json:"director_id,omitempty" db:"director_id"`
	Director    string    `json:"director,omitempty" db:"director"`
	BackdropImg string    `json:"backdrop_img,omitempty" db:"backdrop_img"`
	Rating      float64   `json:"rating" db:"rating"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`

	// Relations
	Genres    []string   `json:"genres"`
	Cast      []string   `json:"cast"`
	Schedules []Schedule `json:"schedules,omitempty"`
}

type MovieRequest struct {
	Title        string          `json:"title" binding:"required"`
	Description  string          `json:"description" binding:"required"`
	ReleaseDate  time.Time       `json:"release_date"`
	Duration     int             `json:"duration" binding:"required"`
	DirectorName string          `json:"director_name" binding:"required"`
	Rating       float64         `json:"rating" binding:"required"`
	PosterImg    string          `json:"poster_img"`
	BackdropImg  string          `json:"backdrop_img"`
	CasterIDs    []int           `json:"caster_ids"`
	GenreIDs     []int           `json:"genre_ids"`
	Schedules    []ScheduleInput `json:"schedules"`
}

type ScheduleInput struct {
	CinemaID int    `json:"cinema_id"`
	TimeID   int    `json:"time_id"`
	Date     string `json:"date"`
	CityID   int    `json:"city_id"`
}
