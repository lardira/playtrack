package player

type RequestCreatePlayer struct {
	Body struct {
		Username string  `json:"username"`
		Img      *string `json:"img" required:"false"`
		Email    *string `json:"email" format:"email" required:"false"`
		Password string  `json:"password"`
	}
}
