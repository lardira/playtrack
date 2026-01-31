package player

import (
	"strings"
	"testing"

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
