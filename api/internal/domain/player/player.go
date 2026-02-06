package player

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/pkg/types"
)

const (
	MinPasswordLength = 8
	MinUsernameLength = 4

	MaxCommentLength = 256
)

const (
	minRating = 1
	maxRating = 100
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
	ErrUsernameMinLen         = fmt.Errorf("username must not be less than %d symbols", MinUsernameLength)
	ErrPasswordMinLen         = fmt.Errorf("password must not be less than %d symbols", MinPasswordLength)
	ErrCompletedEqLessStarted = fmt.Errorf("completed time must not be before or equal to started")
	ErrGameRating             = fmt.Errorf("rating must be in range [%v; %v]", minRating, maxRating)
	ErrCommentLen             = fmt.Errorf("max comment length is %v", MaxCommentLength)
	ErrPlayerID               = fmt.Errorf("invalid player id")
)

var (
	terminatedStatus = map[PlayedGameStatus]struct{}{
		PlayedGameStatusCompleted: struct{}{},
		PlayedGameStatusDropped:   struct{}{},
		PlayedGameStatusRerolled:  struct{}{},
	}

	validPlayedGameStatuses = map[PlayedGameStatus][]PlayedGameStatus{
		PlayedGameStatusAdded:      {PlayedGameStatusInProgress, PlayedGameStatusCompleted, PlayedGameStatusDropped, PlayedGameStatusRerolled},
		PlayedGameStatusInProgress: {PlayedGameStatusCompleted, PlayedGameStatusDropped, PlayedGameStatusRerolled},
		PlayedGameStatusDropped:    {},
		PlayedGameStatusRerolled:   {},
		PlayedGameStatusCompleted:  {},
	}
)

type Player struct {
	ID          string    `json:"id" format:"uuid"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	IsAdmin     bool      `json:"is_admin"`
	Img         *string   `json:"img" required:"false"`
	Email       *string   `json:"email" format:"email" required:"false"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *Player) Valid() error {
	if len(p.Username) < MinUsernameLength {
		return ErrUsernameMinLen
	}
	if len(p.Password) < MinPasswordLength {
		return ErrPasswordMinLen
	}
	return nil
}

type PlayerUpdate struct {
	ID          string
	Username    *string
	Img         *string
	Email       *string
	Password    *string
	Description *string
}

func (p *PlayerUpdate) Valid() error {
	if p.Username != nil && len(*p.Username) < MinUsernameLength {
		return ErrUsernameMinLen
	}
	if p.Password != nil && len(*p.Password) < MinPasswordLength {
		return ErrPasswordMinLen
	}
	return nil
}

type PlayedGame struct {
	ID          int                   `json:"id"`
	PlayerID    string                `json:"player_id"`
	GameID      int                   `json:"game_id"`
	Points      int                   `json:"points"`
	Comment     *string               `json:"comment"`
	Rating      *int                  `json:"rating"`
	Status      PlayedGameStatus      `json:"status"`
	StartedAt   time.Time             `json:"started_at"`
	CompletedAt *time.Time            `json:"completed_at"`
	PlayTime    *types.DurationString `json:"play_time"`
}

type PlayedGameParams struct {
	PlayerID    string
	GameID      int
	Points      int
	Comment     *string
	Rating      *int
	Status      PlayedGameStatus
	StartedAt   time.Time
	CompletedAt *time.Time
	PlayTime    *types.DurationString
}

func NewPlayedGame(params PlayedGameParams) (*PlayedGame, error) {
	nGame := PlayedGame{
		PlayerID: params.PlayerID,
		GameID:   params.GameID,
		Points:   params.Points,
		Status:   PlayedGameStatusAdded,
		PlayTime: params.PlayTime,
	}
	if err := uuid.Validate(nGame.PlayerID); err != nil {
		return nil, ErrPlayerID
	}
	if err := nGame.SetDates(params.StartedAt, params.CompletedAt); err != nil {
		return nil, err
	}
	if params.Rating != nil {
		if err := nGame.SetRating(*params.Rating); err != nil {
			return nil, err
		}
	}
	if params.Comment != nil {
		if err := nGame.SetComment(*params.Comment); err != nil {
			return nil, err
		}
	}
	return &nGame, nil
}

func (pg *PlayedGame) SetDates(startedAt time.Time, completedAt *time.Time) error {
	if completedAt != nil && !startedAt.Before(*completedAt) {
		return ErrCompletedEqLessStarted
	}

	pg.StartedAt = startedAt
	pg.CompletedAt = completedAt
	return nil
}

func (pg *PlayedGame) SetComment(comment string) error {
	if len(comment) > MaxCommentLength {
		return ErrCommentLen
	}
	return nil
}

func (pg *PlayedGame) SetRating(rating int) error {
	if rating < minRating || rating > maxRating {
		return ErrGameRating
	}
	pg.Rating = &rating
	return nil
}

func (pg *PlayedGame) StatusTerminated() bool {
	_, ok := terminatedStatus[pg.Status]
	return ok
}

func (pg *PlayedGame) Complete() error {
	now := time.Now()
	return pg.SetDates(pg.StartedAt, &now)
}

func (pg *PlayedGame) SetStatus(next PlayedGameStatus) error {
	nextMp := validPlayedGameStatuses[pg.Status]
	if ok := slices.Contains(nextMp, next); !ok {
		return fmt.Errorf("next status is not in possible: %v", nextMp)
	}

	pg.Status = next
	return nil
}
