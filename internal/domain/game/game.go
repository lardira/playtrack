package game

import (
	"fmt"
	"net/url"
	"time"
)

const (
	MinGamePoints      = 1
	MinGameHoursToBeat = 1
)

var (
	ErrMinPoints          = fmt.Errorf("game must not have less than %d points", MinGamePoints)
	ErrMinHoursToBeat     = fmt.Errorf("game must not have less than %d hours to beat", MinGameHoursToBeat)
	ErrInvalidGameSiteURL = fmt.Errorf("invalid url")
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
	if g.HoursToBeat < MinGameHoursToBeat {
		return ErrMinHoursToBeat
	}
	if g.Points < MinGamePoints {
		return ErrMinPoints
	}
	if g.URL != nil {
		if _, err := url.ParseRequestURI(*g.URL); err != nil {
			return ErrInvalidGameSiteURL
		}
	}

	return nil
}

func (g *Game) CalculatePoints() {
	// <=2 hours - 1 point
	// >2 hours - each next 4 hours +1 point
	if g.HoursToBeat <= 2 {
		g.Points = 1
	} else {
		// 1 + (n - 2 + d-1) / d
		// e.g. 1 + (10 - 2 + 3) / 4 = 3 points
		g.Points = 1 + (g.HoursToBeat+1)/4
	}
}
