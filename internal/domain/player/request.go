package player

type RequestCreatePlayer struct {
	Body struct {
		Username string  `json:"username"`
		Img      *string `json:"img" required:"false"`
		Email    *string `json:"email" format:"email" required:"false"`
		Password string  `json:"password"`
	}
}

type RequestCreatePlayedGame struct {
	PlayerID string `path:"id" format:"uuid"`
	Body     struct {
		GameID int `json:"game_id"`
	}
}
