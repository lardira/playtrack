package player

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

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
		PlayTime    *DurationString   `json:"play_time" required:"false"`
	}
}

type DurationString struct {
	time.Duration
}

func (d *DurationString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	d.Duration = dur
	return nil
}

func (d DurationString) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d DurationString) Schema(r huma.Registry) *huma.Schema {
	t := reflect.TypeFor[DurationString]()
	r.RegisterTypeAlias(t, reflect.TypeFor[string]())
	return r.Schema(t, true, "")
}
