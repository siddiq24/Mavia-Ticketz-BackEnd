package models

import "time"

type Schedule struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Movie    Movie     `json:"movie"`
	Cinema   Cinema    `json:"cinema"`
	Location Location  `json:"location"`
	Time     Time      `json:"time"`
}

type ScheduleResponse struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Cinema   string    `json:"cinema"`
	Image    string    `json:"image"`
	Time     string    `json:"time"`
	Location string    `json:"location"`
}

type Schedulee struct {
	ID       int       `json:"id" db:"id"`
	MovieID  int       `json:"movie_id" db:"movie_id"`
	CinemaID int       `json:"cinema_id" db:"cinema_id"`
	TimeID   int       `json:"time_id" db:"time_id"`
	Date     time.Time `json:"date" db:"date"`
	CityID   int       `json:"city_id" db:"city_id"`
}

type ScheduleRequest struct {
	CinemaID int       `json:"cinema_id" validate:"required"`
	TimeID   int       `json:"time_id" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	CityID   int       `json:"city_id" validate:"required"`
}

type Cinema struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Location Location
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Time struct {
	ID   int    `json:"id"`
	Time string `json:"time"`
}
