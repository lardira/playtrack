package player

import "time"

type RequestUpdatePlayer struct {
	PlayerID string `path:"id" format:"uuid"`
	Body     struct {
		Username *string `json:"username" required:"false"`
		Img      *string `json:"img" required:"false"`
		Email    *string `json:"email" format:"email" required:"false"`
	}
}

type RequestCreatePlayedGame struct {
	PlayerID string `path:"id" format:"uuid"`
	Body     struct {
		GameID int `json:"game_id"`
	}
}

type RequestUpdatePlayedGame struct {
	PlayerID string `path:"id" format:"uuid"`
	GameID   int    `path:"gameID"`
	Body     struct {
		Points      *int              `json:"points" required:"false"`
		Comment     *string           `json:"comment" required:"false"`
		Rating      *int              `json:"rating" required:"false"`
		Status      *PlayedGameStatus `json:"status" required:"false"`
		CompletedAt *time.Time        `json:"completed_at" required:"false"`
		PlayTime    *time.Duration    `json:"play_time" required:"false"`
	}
}
