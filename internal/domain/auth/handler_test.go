package auth

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	auth "github.com/lardira/playtrack/internal/domain/auth/mock"
	"github.com/lardira/playtrack/internal/domain/player"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
	"github.com/lardira/playtrack/internal/pkg/password"
	"github.com/lardira/playtrack/internal/pkg/testutil"
	"github.com/stretchr/testify/mock"
)

const (
	testSecret = "test"
)

func TestNewHandler(t *testing.T) {
	got := NewHandler(testSecret, auth.NewMockPlayerRepository(t))
	assert.NotEqual(t, nil, got)

	assert.Equal(t, testSecret, got.secret)
	assert.NotEqual(t, nil, got.playerRepository)
}

func TestLogin(t *testing.T) {
	playerRepository := auth.NewMockPlayerRepository(t)
	handler := NewHandler(testSecret, playerRepository)

	playerUsername := "test"
	playerPassword := "test"

	hash, _ := password.Hash(playerPassword)
	testPlayer := player.Player{
		ID:       uuid.NewString(),
		Username: playerUsername,
		Password: hash,
	}

	loginRequest := RequestLoginPlayer{}
	loginRequest.Body.Username = playerUsername
	loginRequest.Body.Password = playerPassword

	playerRepository.
		On("FindOneByUsername", mock.Anything, playerUsername).
		Once().
		Return(&testPlayer, nil)

	resp, err := handler.Login(t.Context(), &loginRequest)
	token := resp.Body.Token

	assert.NoError(t, err)
	assert.NotZero(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		exp, err := t.Claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		if exp.Before(time.Now()) {
			return nil, err
		}
		return []byte(testSecret), nil
	})
	assert.NoError(t, err)
	assert.NotEqual(t, nil, parsedToken)

	sub, err := parsedToken.Claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, testPlayer.ID, sub)
}

func TestRegister(t *testing.T) {
	playerRepository := auth.NewMockPlayerRepository(t)
	handler := NewHandler(testSecret, playerRepository)

	newID := uuid.NewString()
	email := testutil.Faker().Email()
	req := RequestRegisterCreatePlayer{}
	req.Body.Username = testutil.Faker().Username()
	req.Body.Password = testutil.Faker().Password(true, true, true, true, false, player.MinPasswordLength)
	req.Body.Email = &email

	var constructedPlayer *player.Player

	playerRepository.
		On(
			"Insert",
			mock.Anything,
			mock.MatchedBy(func(p *player.Player) bool {
				if p == nil || p.Username == "" || p.Password == "" {
					return false
				}
				constructedPlayer = p
				return true
			})).
		Once().
		Return(newID, nil)

	resp, err := handler.RegisterPlayer(t.Context(), &req)
	assert.NoError(t, err)
	assert.Equal(t, newID, resp.Body.ID)

	assert.True(t, password.CompareHash(req.Body.Password, constructedPlayer.Password))
}

func TestSetPassword(t *testing.T) {
	playerRepository := auth.NewMockPlayerRepository(t)
	handler := NewHandler(testSecret, playerRepository)

	newID := uuid.NewString()
	ctx := ctxutil.SetPlayerID(t.Context(), newID)

	req := RequestSetPassword{}
	req.Body.Password = testutil.Faker().Password(true, true, true, true, false, player.MinPasswordLength)

	var constructedPlayer *player.PlayerUpdate

	playerRepository.
		On(
			"Update",
			ctx,
			mock.MatchedBy(func(p *player.PlayerUpdate) bool {
				if p == nil || p.Password == nil || p.ID == "" {
					return false
				}
				constructedPlayer = p
				return true
			})).
		Once().
		Return(newID, nil)

	resp, err := handler.SetPassword(ctx, &req)
	assert.NoError(t, err)
	assert.Equal(t, newID, resp.Body.ID)

	assert.True(t, password.CompareHash(req.Body.Password, *constructedPlayer.Password))
}

func TestIssueToken(t *testing.T) {
	handler := NewHandler(testSecret, auth.NewMockPlayerRepository(t))
	playerID := uuid.NewString()

	token, err := handler.issueToken(playerID)
	assert.NoError(t, err)
	assert.NotZero(t, token)
}
