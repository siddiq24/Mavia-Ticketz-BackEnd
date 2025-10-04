package models

type TicketSalesByMovie struct {
	MovieTitle string `json:"movie_title"`
	Week       int    `json:"week"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	TotalSales int    `json:"total_sales"`
}

type TicketSalesByCategoryLocation struct {
	Category   string `json:"category"`
	Location   string `json:"location"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	TotalSales int    `json:"total_sales"`
}
