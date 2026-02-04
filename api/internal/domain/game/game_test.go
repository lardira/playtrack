package game

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/lardira/playtrack/internal/pkg/testutil"
)

func TestValidGame(t *testing.T) {
	url := testutil.Faker().URL()
	invalidURL := "example.cra"

	tcases := []struct {
		name string
		game Game
		want error
	}{
		{
			"valid game",
			Game{
				HoursToBeat: MinGameHoursToBeat,
				URL:         &url,
				Points:      MinGamePoints,
			},
			nil,
		},
		{
			"invalid hours to beat",
			Game{
				HoursToBeat: MinGameHoursToBeat - 1,
				URL:         &url,
				Points:      MinGamePoints,
			},
			ErrMinHoursToBeat,
		},
		{
			"invalid url",
			Game{
				HoursToBeat: MinGameHoursToBeat,
				URL:         &invalidURL,
				Points:      MinGamePoints,
			},
			ErrInvalidGameSiteURL,
		},
		{
			"invalid points",
			Game{
				HoursToBeat: MinGameHoursToBeat,
				URL:         &url,
				Points:      MinGamePoints - 1,
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
		{"min game hours points", MinGameHoursToBeat, 1},
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
