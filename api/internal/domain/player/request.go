package player

import (
	"time"

	"github.com/lardira/playtrack/internal/pkg/types"
)

type RequestUpdatePlayer struct {
	PlayerID string `path:"id" format:"uuid"`
	Body     struct {
		Username    *string `json:"username" minLength:"4" maxLength:"32" required:"false"`
		Img         *string `json:"img" format:"uri" required:"false"`
		Email       *string `json:"email" format:"email" required:"false"`
		Description *string `json:"description" maxLength:"128" required:"false"`
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
		Points      *int                  `json:"points" required:"false"`
		Comment     *string               `json:"comment" required:"false"`
		Rating      *int                  `json:"rating" required:"false"`
		Status      *PlayedGameStatus     `json:"status" required:"false"`
		StartedAt   *time.Time            `json:"started_at" dependentRequired:"completed_at" required:"false"`
		CompletedAt *time.Time            `json:"completed_at" dependentRequired:"started_at" required:"false"`
		PlayTime    *types.DurationString `json:"play_time" required:"false"`
	}
}
