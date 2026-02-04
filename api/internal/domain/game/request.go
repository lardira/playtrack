package game

type RequestCreateGame struct {
	Body struct {
		HoursToBeat int     `json:"hours_to_beat"`
		Title       string  `json:"title"`
		URL         *string `json:"url" required:"false" format:"uri"`
	}
}
