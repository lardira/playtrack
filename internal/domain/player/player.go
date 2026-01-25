package player

import (
	"fmt"
	"slices"
	"time"
)

const (
	minPasswordLength = 8
	minUsernameLength = 4
)

type PlayedGameStatus string

const (
	PlayedGameStatusAdded      PlayedGameStatus = "added"
	PlayedGameStatusInProgress PlayedGameStatus = "in_progress"
	PlayedGameStatusCompleted  PlayedGameStatus = "completed"
	PlayedGameStatusDropped    PlayedGameStatus = "dropped"
	PlayedGameStatusRerolled   PlayedGameStatus = "rerolled"
)

var (
	terminatedStatus = []PlayedGameStatus{
		PlayedGameStatusCompleted,
		PlayedGameStatusDropped,
		PlayedGameStatusRerolled,
	}
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

type PlayedGame struct {
	ID          int              `json:"id"`
	PlayerID    string           `json:"player_id"`
	GameID      int              `json:"game_id"`
	Points      int              `json:"points"`
	Comment     *string          `json:"comment"`
	Rating      *int             `json:"rating"`
	Status      PlayedGameStatus `json:"status"`
	StartedAt   time.Time        `json:"started_at"`
	CompletedAt *time.Time       `json:"completed_at"`
	PlayTime    *time.Duration   `json:"play_time"`
}

func (pg *PlayedGame) Valid() error {
	// TODO: validation
	return nil
}

func (pg *PlayedGame) StatusTerminated() bool {
	return slices.Contains(terminatedStatus, pg.Status)
}

type LeaderboardPlayer struct {
	PlayerID  string `json:"player_id"`
	Completed int    `json:"completed"`
	Total     int    `json:"total"`
	Dropped   int    `json:"dropped"`
	Rerolled  int    `json:"rerolled"`
}
