package game

import (
	"fmt"
	"net/url"
	"time"
)

const (
	minGamePoints      = 1
	minGameHoursToBeat = 2
)

var (
	ErrMinPoints          = fmt.Errorf("game must not have less than %d points", minGamePoints)
	ErrMinHoursToBeat     = fmt.Errorf("game must not have less than %d hours to beat", minGameHoursToBeat)
	ErrInvalidGamesiteURL = fmt.Errorf("invalid url")
)

type Game struct {
	ID          int       `json:"id"`
	Points      int       `json:"points"`
	HoursToBeat int       `json:"hours_to_beat"`
	Title       string    `json:"title"`
	URL         *string   `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

func (g *Game) Valid() error {
	if g.HoursToBeat < minGameHoursToBeat {
		return ErrMinHoursToBeat
	}
	if g.Points < minGamePoints {
		return ErrMinPoints
	}
	if g.URL != nil {
		if _, err := url.ParseRequestURI(*g.URL); err != nil {
			return ErrInvalidGamesiteURL
		}
	}

	return nil
}
