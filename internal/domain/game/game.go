package game

type Game struct {
	ID     int    `json:"id"`
	Points int    `json:"points"`
	Title  string `json:"title"`
}
