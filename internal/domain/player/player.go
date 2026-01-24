package player

import (
	"fmt"
	"time"
)

const (
	minPasswordLength = 8
	minUsernameLength = 4
)

type GamePlayedStatus string

const (
	GamePlayedStatusAdded      GamePlayedStatus = "added"
	GamePlayedStatusInProgress GamePlayedStatus = "in_progress"
	GamePlayedStatusCompleted  GamePlayedStatus = "completed"
	GamePlayedStatusDropped    GamePlayedStatus = "dropped"
	GamePlayedStatusRerolled   GamePlayedStatus = "rerolled"
)

type Player struct {
	ID        string    `json:"id" format:"uuid"`
	Username  string    `json:"username"`
	Img       *string   `json:"img" required:"false"`
	Email     *string   `json:"email" format:"email" required:"false"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Player) Valid() error {

	if len(p.Username) < minUsernameLength {
		return fmt.Errorf("username must not be less than %d symbols", minUsernameLength)
	}

	if len(p.Password) < minPasswordLength {
		return fmt.Errorf("password must not be less than %d symbols", minPasswordLength)
	}

	return nil
}

type GamePlayed struct {
	ID          int              `json:"id"`
	GameID      int              `json:"game_id"`
	Scores      int              `json:"scores"`
	PlayerID    string           `json:"player_id"`
	Status      GamePlayedStatus `json:"status"`
	Comment     *string          `json:"comment"`
	Rating      *int             `json:"rating"`
	StartedAt   time.Time        `json:"started_at"`
	CompletedAt *time.Time       `json:"completed_at" required:"false"`
}

type LeaderboardPlayer struct {
	PlayerID  string `json:"player_id"`
	Completed int    `json:"completed"`
	Total     int    `json:"total"`
	Dropped   int    `json:"dropped"`
	Rerolled  int    `json:"rerolled"`
}
