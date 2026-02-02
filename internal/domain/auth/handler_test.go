package auth

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	got := NewHandler(testSecret, NewMockPlayerRepository(t))
	assert.NotEqual(t, nil, got)

	assert.Equal(t, testSecret, got.secret)
	assert.NotEqual(t, nil, got.playerRepository)
}

func TestLogin(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)
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
	playerRepository := NewMockPlayerRepository(t)
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
	playerRepository := NewMockPlayerRepository(t)
	handler := NewHandler(testSecret, playerRepository)

	playerID := uuid.NewString()
	username := testutil.Faker().Username()
	ctx := ctxutil.SetPlayer(t.Context(), ctxutil.CtxPlayer{ID: playerID})

	req := RequestSetPassword{}
	req.Body.Username = username
	req.Body.Password = testutil.Faker().Password(true, true, true, true, false, player.MinPasswordLength)

	var constructedPlayer *player.PlayerUpdate

	playerRepository.
		On("FindOneByUsername", ctx, username).
		Once().
		Return(&player.Player{ID: playerID, Username: username}, nil)

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
		Return(playerID, nil)

	resp, err := handler.SetPassword(ctx, &req)
	assert.NoError(t, err)
	assert.Equal(t, playerID, resp.Body.ID)

	assert.True(t, password.CompareHash(req.Body.Password, *constructedPlayer.Password))
}

func TestSetPassword_DifferentPlayer(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)
	handler := NewHandler(testSecret, playerRepository)

	newID := uuid.NewString()
	diffID := uuid.NewString()
	username := testutil.Faker().Username()
	ctx := ctxutil.SetPlayer(t.Context(), ctxutil.CtxPlayer{ID: newID})

	req := RequestSetPassword{}
	req.Body.Username = username
	req.Body.Password = testutil.Faker().Password(true, true, true, true, false, player.MinPasswordLength)

	playerRepository.
		On("FindOneByUsername", ctx, username).
		Once().
		Return(&player.Player{ID: diffID, Username: username}, nil)

	playerRepository.AssertNotCalled(t, "Update")

	_, err := handler.SetPassword(ctx, &req)
	assert.Error(t, err)
}

func TestSetPassword_DifferentPlayer_AsAdmin(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)
	handler := NewHandler(testSecret, playerRepository)

	adminID := uuid.NewString()
	diffID := uuid.NewString()
	username := testutil.Faker().Username()
	ctx := ctxutil.SetPlayer(t.Context(), ctxutil.CtxPlayer{ID: adminID, IsAdmin: true})

	req := RequestSetPassword{}
	req.Body.Username = username
	req.Body.Password = testutil.Faker().Password(true, true, true, true, false, player.MinPasswordLength)

	var constructedPlayer *player.PlayerUpdate

	playerRepository.
		On("FindOneByUsername", ctx, username).
		Once().
		Return(&player.Player{ID: diffID, Username: username}, nil)

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
		Return(diffID, nil)

	resp, err := handler.SetPassword(ctx, &req)
	assert.NoError(t, err)
	assert.Equal(t, diffID, resp.Body.ID)

	assert.Equal(t, diffID, constructedPlayer.ID)

	assert.True(t, password.CompareHash(req.Body.Password, *constructedPlayer.Password))
}

func TestIssueToken(t *testing.T) {
	var p player.Player
	testutil.Faker().Struct(&p)

	handler := NewHandler(testSecret, NewMockPlayerRepository(t))

	token, err := handler.issueToken(&p)
	assert.NoError(t, err)
	assert.NotZero(t, token)
}
