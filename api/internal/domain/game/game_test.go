package game

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/lardira/playtrack/internal/pkg/testutil"
)

func TestValidGame(t *testing.T) {
	url := testutil.Faker().URL()
	title := testutil.Faker().Word()
	validHours := testutil.Faker().IntRange(MinGameHoursToBeat, 100)
	invalidURL := "example.cra"

	tcases := []struct {
		name        string
		gameTitle   string
		hoursToBeat int
		url         *string
		want        error
	}{
		{
			"valid game",
			title,
			validHours,
			&url,
			nil,
		},
		{
			"invalid hours to beat",
			title,
			MinGameHoursToBeat - 1,
			&url,
			ErrMinHoursToBeat,
		},
		{
			"invalid url",
			title,
			validHours,
			&invalidURL,
			ErrInvalidGameSiteURL,
		},
		{
			"invalid title",
			strings.Repeat("t", MinGameTitleLen-1),
			validHours,
			&url,
			ErrMinTitleLen,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGame(tt.gameTitle, tt.hoursToBeat, tt.url)

			assert.IsError(t, err, tt.want)
			if err != nil {
				return
			}
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
