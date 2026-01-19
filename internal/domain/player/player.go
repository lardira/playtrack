package player

type Player struct {
	ID       string  `json:"id" format:"uuid"`
	Username string  `json:"username"`
	Img      string  `json:"img"`
	Email    *string `json:"email" format:"email" required:"false"`
	Password string  `json:"password"`
}
