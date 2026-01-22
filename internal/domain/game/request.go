package game

type RequestCreate struct {
	Body struct {
		Points      int     `json:"points"`
		HoursToBeat int     `json:"hours_to_beat"`
		Title       string  `json:"title"`
		URL         *string `json:"url" required:"false" format:"uri"`
	}
}
