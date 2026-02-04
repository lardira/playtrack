package game

type RequestCreateGame struct {
	Body struct {
		HoursToBeat int     `json:"hours_to_beat" minimum:"1"`
		Title       string  `json:"title" minLength:"2"`
		URL         *string `json:"url" required:"false" format:"uri"`
	}
}
