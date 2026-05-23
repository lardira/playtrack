package game

import (
	"fmt"
	"net/url"
	"time"
)

const (
	MinGamePoints      = 1
	MinGameHoursToBeat = 1
	MinGameTitleLen    = 3
)

var (
	ErrMinHoursToBeat     = fmt.Errorf("game must not have less than %d hours to beat", MinGameHoursToBeat)
	ErrInvalidGameSiteURL = fmt.Errorf("invalid url")
	ErrMinTitleLen        = fmt.Errorf("title length must not be less than %d", MinGameTitleLen)
)

type Game struct {
	ID          int       `json:"id"`
	Points      int       `json:"points"`
	HoursToBeat int       `json:"hours_to_beat"`
	Title       string    `json:"title"`
	URL         *string   `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewGame(title string, hoursToBeat int, url *string) (*Game, error) {
	nGame := Game{
		Title:     title,
		CreatedAt: time.Now(),
	}
	if url != nil {
		if err := nGame.SetURL(*url); err != nil {
			return nil, err
		}
	}
	if err := nGame.SetHoursToBeat(hoursToBeat); err != nil {
		return nil, err
	}
	if len(nGame.Title) < MinGameTitleLen {
		return nil, ErrMinTitleLen
	}
	return &nGame, nil
}

// CalculatePoints calculates and sets points in game
//
// Rules:
//   - less than 2 hours = 1 point
//   - more than 2 hours = each next 4 hours +1 point
func (g *Game) CalculatePoints() {
	if g.HoursToBeat <= 2 {
		g.Points = 1
	} else {
		// points = 1 + (n - 2 + d-1) / d
		// e.g. 1 + (10 - 2 + 3) / 4 = 3 points
		g.Points = 1 + (g.HoursToBeat+1)/4
	}
}

func (g *Game) SetHoursToBeat(hours int) error {
	if hours < MinGameHoursToBeat {
		return ErrMinHoursToBeat
	}

	g.HoursToBeat = hours
	g.CalculatePoints()
	return nil
}

func (g *Game) SetURL(u string) error {
	if _, err := url.ParseRequestURI(u); err != nil {
		return ErrInvalidGameSiteURL
	}
	g.URL = &u

	return nil
}
