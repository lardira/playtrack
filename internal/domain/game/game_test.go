package game

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

var (
	validURL   = "http://example.test.test"
	invalidURL = "example.cra"
)

func TestValidGame(t *testing.T) {
	tcases := []struct {
		name string
		game Game
		want error
	}{
		{
			"valid game",
			Game{
				HoursToBeat: minGameHoursToBeat,
				URL:         &validURL,
				Points:      minGamePoints,
			},
			nil,
		},
		{
			"invalid hours to beat",
			Game{
				HoursToBeat: minGameHoursToBeat - 1,
				URL:         &validURL,
				Points:      minGamePoints,
			},
			ErrMinHoursToBeat,
		},
		{
			"invalid url",
			Game{
				HoursToBeat: minGameHoursToBeat,
				URL:         &invalidURL,
				Points:      minGamePoints,
			},
			ErrInvalidGamesiteURL,
		},
		{
			"invalid points",
			Game{
				HoursToBeat: minGameHoursToBeat,
				URL:         &validURL,
				Points:      minGamePoints - 1,
			},
			ErrMinPoints,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.game.Valid()

			assert.IsError(t, err, tt.want)
		})
	}
}

func TestGameCalculatePoints(t *testing.T) {
	tcases := []struct {
		name  string
		hours int
		want  int
	}{
		{"min game hours points", minGameHoursToBeat, 1},
		{"2 hours 1 point", 2, 1},
		{"9 hours 3 points", 9, 3},
		{"10 hours 3 points", 10, 3},
		{"11 hours 4 points", 11, 4},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			game := Game{HoursToBeat: tt.hours}

			game.CalculatePoints()

			assert.Equal(t, tt.want, game.Points)
		})
	}
}
