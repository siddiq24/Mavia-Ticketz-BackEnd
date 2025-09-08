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
