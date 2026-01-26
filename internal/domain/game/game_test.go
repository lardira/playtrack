package game

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

var (
	validURL   = "http://example.test.test"
	invalidURL = "example.cra"
)

func TestValidErr(t *testing.T) {
	tcases := map[string]struct {
		game Game
		want error
	}{
		"invalid hours to beat": {
			Game{
				HoursToBeat: minGameHoursToBeat - 1,
				URL:         &validURL,
				Points:      minGamePoints,
			},
			ErrMinHoursToBeat,
		},
		"invalid url": {
			Game{
				HoursToBeat: minGameHoursToBeat,
				URL:         &invalidURL,
				Points:      minGamePoints,
			},
			ErrInvalidGamesiteURL,
		},
		"invalid points": {
			Game{
				HoursToBeat: minGameHoursToBeat,
				URL:         &validURL,
				Points:      minGamePoints - 1,
			},
			ErrMinPoints,
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tcase.game.Valid()

			assert.IsError(t, err, tcase.want)
		})
	}
}

func TestGameCalculatePoints(t *testing.T) {
	tcases := map[string]struct {
		game Game
		want int
	}{
		"min game hours points": {
			Game{HoursToBeat: minGameHoursToBeat},
			1,
		},
		"2 hours 1 point": {
			Game{HoursToBeat: 2},
			1,
		},
		"9 hours 3 points": {
			Game{HoursToBeat: 9},
			3,
		},
		"10 hours 3 points": {
			Game{HoursToBeat: 10},
			3,
		},
		"11 hours 4 points": {
			Game{HoursToBeat: 11},
			4,
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tcase.game.CalculatePoints()

			assert.Equal(t, tcase.want, tcase.game.Points)
		})
	}
}
