package player

import (
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/pkg/testutil"
	"github.com/lardira/playtrack/internal/pkg/types"
)

func TestPlayerValid(t *testing.T) {
	url := testutil.Faker().URL()
	email := testutil.Faker().Email()
	password := testutil.Faker().Password(true, true, true, true, false, MinPasswordLength)

	validParams := PlayerParams{
		Username:    testutil.Faker().Username(),
		Img:         &url,
		Email:       &email,
		Password:    password,
		Description: &testutil.Faker().Address().Address,
	}

	tcases := []struct {
		name   string
		params func() PlayerParams
		err    error
	}{
		{
			"valid",
			func() PlayerParams {
				return validParams
			},
			nil,
		},
		{
			"password less than min",
			func() PlayerParams {
				p := validParams
				p.Password = strings.Repeat("a", MinPasswordLength-1)
				return p
			},
			ErrPasswordMinLen,
		},
		{
			"username less than min",
			func() PlayerParams {
				p := validParams
				p.Username = strings.Repeat("a", MinUsernameLength-1)
				return p
			},
			ErrUsernameMinLen,
		},
		{
			"invalid image url",
			func() PlayerParams {
				p := validParams
				img := "test string"
				p.Img = &img
				return p
			},
			ErrInvalidURL,
		},
		{
			"invalid email address",
			func() PlayerParams {
				p := validParams
				email := testutil.Faker().URL()
				p.Email = &email
				return p
			},
			ErrInvalidEmail,
		},
		{
			"invalid description length",
			func() PlayerParams {
				p := validParams
				descr := strings.Repeat("t", MaxTextLength+1)
				p.Description = &descr
				return p
			},
			ErrTextFieldLen,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.params()
			_, err := NewPlayer(params)
			assert.IsError(t, tt.err, err)
		})
	}
}

func TestPlayedGame(t *testing.T) {
	now := time.Now()
	comment := testutil.Faker().Sentence()
	rating := testutil.Faker().IntRange(minRating, maxRating)
	validParams := PlayedGameParams{
		PlayerID:  uuid.NewString(),
		GameID:    testutil.Faker().Int(),
		Points:    0,
		Comment:   &comment,
		Rating:    &rating,
		Status:    PlayedGameStatusAdded,
		StartedAt: now,
		PlayTime:  &types.DurationString{},
	}

	tcases := []struct {
		name    string
		params  func() PlayedGameParams
		wantErr error
	}{
		{
			"valid",
			func() PlayedGameParams {
				return validParams
			},
			nil,
		},
		{
			"completed before started",
			func() PlayedGameParams {
				p := validParams
				p.CompletedAt = &now
				return p
			},
			ErrCompletedEqLessStarted,
		},
		{
			"min rating",
			func() PlayedGameParams {
				p := validParams
				r := minRating - 1
				p.Rating = &r
				return p
			},
			ErrGameRating,
		},
		{
			"max rating",
			func() PlayedGameParams {
				p := validParams
				r := maxRating + 1
				p.Rating = &r
				return p
			},
			ErrGameRating,
		},
		{
			"invalid player id",
			func() PlayedGameParams {
				p := validParams
				p.PlayerID = testutil.Faker().Word()
				return p
			},
			ErrPlayerID,
		},
		{
			"comment too long",
			func() PlayedGameParams {
				p := validParams
				c := strings.Repeat("t", MaxTextLength+1)
				p.Comment = &c
				return p
			},
			ErrTextFieldLen,
		},
	}

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPlayedGame(tt.params())
			assert.IsError(t, err, tt.wantErr)
		})
	}
}

func TestPlayedGameSetStatus(t *testing.T) {
	tcases := []struct {
		name      string
		current   PlayedGameStatus
		next      PlayedGameStatus
		expectErr bool
	}{
		{
			name:      "valid next in progress",
			current:   PlayedGameStatusAdded,
			next:      PlayedGameStatusInProgress,
			expectErr: false,
		},
		{
			name:      "valid next in completed",
			current:   PlayedGameStatusAdded,
			next:      PlayedGameStatusCompleted,
			expectErr: false,
		},
		{
			name:      "valid in progresss next in completed",
			current:   PlayedGameStatusInProgress,
			next:      PlayedGameStatusCompleted,
			expectErr: false,
		},
		{
			name:      "invalid dropped next rerolled",
			current:   PlayedGameStatusDropped,
			next:      PlayedGameStatusRerolled,
			expectErr: true,
		},
		{
			name:      "invalid completed in progress",
			current:   PlayedGameStatusCompleted,
			next:      PlayedGameStatusInProgress,
			expectErr: true,
		},
	}
	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			pg := validPlayedGame()
			pg.Status = tt.current

			err := pg.SetStatus(tt.next)
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}

func TestPlayedGameStatusTerminated(t *testing.T) {
	tcases := []struct {
		name   string
		status PlayedGameStatus
		ok     bool
	}{
		{
			name:   "terminated completed",
			status: PlayedGameStatusCompleted,
			ok:     true,
		},
		{
			name:   "terminated dropped",
			status: PlayedGameStatusDropped,
			ok:     true,
		},
		{
			name:   "terminated in progress",
			status: PlayedGameStatusInProgress,
			ok:     false,
		},
	}
	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			pg := validPlayedGame()
			pg.Status = tt.status

			assert.Equal(t, tt.ok, pg.StatusTerminated())
		})
	}
}
func TestPlayedGameComplete(t *testing.T) {
	pg := validPlayedGame()
	pg.StartedAt = time.Now().Add(-2 * time.Hour)

	err := pg.Complete()
	assert.NoError(t, err)
	assert.NotEqual(t, nil, pg.CompletedAt)
}

func validPlayer() Player {
	url := testutil.Faker().URL()
	email := testutil.Faker().Email()
	password := testutil.Faker().Password(true, true, true, true, false, MinPasswordLength)

	p, _ := NewPlayer(PlayerParams{
		Username:    testutil.Faker().Username(),
		Img:         &url,
		Email:       &email,
		Password:    password,
		Description: &testutil.Faker().Address().Address,
	})
	return *p
}

func validPlayedGame() PlayedGame {
	now := time.Now()
	comment := testutil.Faker().Sentence()
	rating := testutil.Faker().IntRange(minRating, maxRating)

	pg, _ := NewPlayedGame(PlayedGameParams{
		PlayerID:  uuid.NewString(),
		GameID:    testutil.Faker().Int(),
		Points:    0,
		Comment:   &comment,
		Rating:    &rating,
		Status:    PlayedGameStatusAdded,
		StartedAt: now,
		PlayTime:  &types.DurationString{},
	})
	return *pg
}
