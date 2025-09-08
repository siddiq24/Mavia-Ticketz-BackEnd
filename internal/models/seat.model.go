package models

type Seat struct {
	ID     int    `json:"id"`
	Cols   string `json:"cols"`
	Rows   int    `json:"rows"`
	Status string `json:"status"`
}

type AvailableSeat struct {
	ID     int    `json:"id"`
	Cols   string `json:"cols"`
	Rows   int    `json:"rows"`
	Status string `json:"status"`
}
