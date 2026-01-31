package player

import (
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/pkg/testutil"
)

func TestPlayerValid(t *testing.T) {
	tcases := []struct {
		name   string
		player func() Player
		err    error
	}{
		{
			"valid",
			validPlayer,
			nil,
		},
		{
			"password less than min",
			func() Player {
				p := validPlayer()
				p.Password = strings.Repeat("a", MinPasswordLength-1)
				return p
			},
			ErrPasswordMinLen,
		},
		{
			"username less than min",
			func() Player {
				p := validPlayer()
				p.Username = strings.Repeat("a", MinUsernameLength-1)
				return p
			},
			ErrUsernameMinLen,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.player()
			gotErr := p.Valid()

			assert.IsError(t, tt.err, gotErr)
		})
	}
}

func TestPlayerUpdateValid(t *testing.T) {
	tcases := []struct {
		name   string
		player func() PlayerUpdate
		err    error
	}{
		{
			"valid",
			validPlayerUpdate,
			nil,
		},
		{
			"password less than min",
			func() PlayerUpdate {
				p := validPlayerUpdate()
				newPass := strings.Repeat("a", MinPasswordLength-1)
				p.Password = &newPass
				return p
			},
			ErrPasswordMinLen,
		},
		{
			"username less than min",
			func() PlayerUpdate {
				p := validPlayerUpdate()
				newUsername := strings.Repeat("a", MinUsernameLength-1)
				p.Username = &newUsername
				return p
			},
			ErrUsernameMinLen,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.player()
			gotErr := p.Valid()

			assert.IsError(t, tt.err, gotErr)
		})
	}
}

func TestPlayedGameValid(t *testing.T) {
	tcases := []struct {
		name string
		game func() PlayedGame
		err  error
	}{
		{
			"valid",
			validPlayedGame,
			nil,
		},
		{
			"completed before started",
			func() PlayedGame {
				g := validPlayedGame()
				g.StartedAt = time.Now()
				return g
			},
			ErrCompletedBeforeStarted,
		},
		{
			"min rating",
			func() PlayedGame {
				g := validPlayedGame()
				r := minRating - 1
				g.Rating = &r
				return g
			},
			ErrGameRating,
		},
		{
			"max rating",
			func() PlayedGame {
				g := validPlayedGame()
				r := maxRating + 1
				g.Rating = &r
				return g
			},
			ErrGameRating,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.game()
			gotErr := g.Valid()

			assert.IsError(t, tt.err, gotErr)
		})
	}
}

func TestStatusTerminated(t *testing.T) {
	tcases := []struct {
		name string
		game func() PlayedGame
		ok   bool
	}{
		{
			"completed",
			func() PlayedGame {
				g := validPlayedGame()
				g.Status = PlayedGameStatusCompleted
				return g
			},
			true,
		},
		{
			"in progress",
			func() PlayedGame {
				g := validPlayedGame()
				g.Status = PlayedGameStatusInProgress
				return g
			},
			false,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.game()
			ok := g.StatusTerminated()

			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestStatusNextValid(t *testing.T) {
	game := validPlayedGame()
	game.Status = PlayedGameStatusInProgress

	err := game.StatusNextValid(PlayedGameStatusCompleted)
	assert.NoError(t, err)

	game.Status = PlayedGameStatusCompleted

	err = game.StatusNextValid(PlayedGameStatusInProgress)
	assert.Error(t, err)
}

func TestPlayedGameUpdateValid(t *testing.T) {
	tcases := []struct {
		name string
		game func() PlayedGameUpdate
		err  error
	}{
		{
			"valid",
			validPlayedGameUpdate,
			nil,
		},
		{
			"min rating",
			func() PlayedGameUpdate {
				g := validPlayedGameUpdate()
				r := minRating - 1
				g.Rating = &r
				return g
			},
			ErrGameRating,
		},
		{
			"max rating",
			func() PlayedGameUpdate {
				g := validPlayedGameUpdate()
				r := maxRating + 1
				g.Rating = &r
				return g
			},
			ErrGameRating,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.game()
			err := g.Valid()

			assert.IsError(t, tt.err, err)
		})
	}
}

func validPlayer() Player {
	url := testutil.Faker().URL()
	email := testutil.Faker().Email()
	password := testutil.Faker().Password(true, true, true, true, false, MinPasswordLength)

	return Player{
		ID:        uuid.NewString(),
		Username:  testutil.Faker().Username(),
		Img:       &url,
		Email:     &email,
		Password:  password,
		CreatedAt: testutil.Faker().Date(),
	}
}

func validPlayedGame() PlayedGame {
	comment := testutil.Faker().Comment()
	completedAt := time.Now()
	playTime := 1 * time.Hour
	rating := testutil.Faker().IntRange(minRating, maxRating)

	return PlayedGame{
		ID:          testutil.Faker().Int(),
		PlayerID:    uuid.NewString(),
		GameID:      testutil.Faker().Int(),
		Points:      testutil.Faker().Int(),
		Comment:     &comment,
		Rating:      &rating,
		Status:      PlayedGameStatusCompleted,
		StartedAt:   time.Now().Add(-1 * time.Hour),
		CompletedAt: &completedAt,
		PlayTime:    &playTime,
	}
}

func validPlayerUpdate() PlayerUpdate {
	url := testutil.Faker().URL()
	username := testutil.Faker().Username()
	email := testutil.Faker().Email()
	password := testutil.Faker().Password(true, true, true, true, false, MinPasswordLength)

	return PlayerUpdate{
		ID:       uuid.NewString(),
		Img:      &url,
		Username: &username,
		Email:    &email,
		Password: &password,
	}
}

func validPlayedGameUpdate() PlayedGameUpdate {
	comment := testutil.Faker().Comment()
	completedAt := time.Now()
	playTime := 1 * time.Hour
	rating := testutil.Faker().IntRange(minRating, maxRating)
	points := testutil.Faker().Int()
	status := PlayedGameStatusCompleted

	return PlayedGameUpdate{
		ID:          testutil.Faker().Int(),
		Points:      &points,
		Comment:     &comment,
		Rating:      &rating,
		Status:      &status,
		CompletedAt: &completedAt,
		PlayTime:    &playTime,
	}
}
