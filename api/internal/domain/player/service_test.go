package player

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/domain/game"
	"github.com/lardira/playtrack/internal/pkg/testutil"
	"github.com/lardira/playtrack/internal/pkg/types"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)

	service := NewService(
		playerRepository,
		NewMockPlayedGameRepository(t),
		NewMockGameService(t),
	)

	playerRepository.
		On("FindAll", t.Context()).
		Once().
		Return([]Player{}, nil)

	players, err := service.GetAll(t.Context())
	assert.NoError(t, err)
	assert.NotEqual(t, nil, players)
}

func TestGetOne(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)

	service := NewService(
		playerRepository,
		NewMockPlayedGameRepository(t),
		NewMockGameService(t),
	)

	playerRepository.
		On("FindOne", t.Context(), mock.Anything).
		Once().
		Return(&Player{}, nil)

	players, err := service.GetOne(t.Context(), "some-id")
	assert.NoError(t, err)
	assert.NotEqual(t, nil, players)
}

func TestGetOneByUsername(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)

	service := NewService(
		playerRepository,
		NewMockPlayedGameRepository(t),
		NewMockGameService(t),
	)

	playerRepository.
		On("FindOneByUsername", t.Context(), mock.Anything).
		Once().
		Return(&Player{}, nil)

	players, err := service.GetOneByUsername(t.Context(), "some-username")
	assert.NoError(t, err)
	assert.NotEqual(t, nil, players)
}

func TestCreate(t *testing.T) {
	newID := uuid.NewString()

	playerRepository := NewMockPlayerRepository(t)

	service := NewService(
		playerRepository,
		NewMockPlayedGameRepository(t),
		NewMockGameService(t),
	)

	playerRepository.
		On("Insert", t.Context(), mock.AnythingOfType("*player.Player")).
		Once().
		Return(newID, nil)

	url := testutil.Faker().URL()
	email := testutil.Faker().Email()
	password := testutil.Faker().Password(true, true, true, true, false, MinPasswordLength)

	playerID, err := service.Create(t.Context(), PlayerParams{
		Username:    testutil.Faker().Username(),
		Img:         &url,
		Email:       &email,
		Password:    password,
		Description: &testutil.Faker().Address().Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, newID, playerID)
}

func TestUpdate(t *testing.T) {
	validPlayer := validPlayer()
	newUsername := testutil.Faker().Username()
	newPassword := testutil.Faker().Password(true, true, true, true, false, MinPasswordLength)
	url := testutil.Faker().URL()
	email := testutil.Faker().Email()
	description := testutil.Faker().Blurb()

	playerRepository := NewMockPlayerRepository(t)

	service := NewService(
		playerRepository,
		NewMockPlayedGameRepository(t),
		NewMockGameService(t),
	)

	playerRepository.
		On("FindOne", t.Context(), mock.Anything).
		Once().
		Return(&validPlayer, nil)

	var upd Player
	playerRepository.
		On("Update", t.Context(), mock.MatchedBy(func(p *Player) bool {
			if p.ID != validPlayer.ID {
				return false
			}
			upd = *p
			return true
		})).
		Once().
		Return(validPlayer.ID, nil)

	playerID, err := service.Update(t.Context(), validPlayer.ID, PlayerUpdate{
		Username:    &newUsername,
		Password:    &newPassword,
		Img:         &url,
		Email:       &email,
		Description: &description,
	})
	assert.NoError(t, err)

	assert.Equal(t, playerID, validPlayer.ID)
	assert.Equal(t, url, *upd.Img)
	assert.Equal(t, email, *upd.Email)
	assert.Equal(t, description, *upd.Description)

	assert.Equal(t, newUsername, upd.Username)
	assert.NotEqual(t, newPassword, upd.Password)
	assert.True(t, upd.CheckPassword(newPassword))
}

func TestGetAllPlayedGames(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	service := NewService(
		playerRepository,
		playedGameRepository,
		NewMockGameService(t),
	)

	playedGameRepository.
		On("FindAll", t.Context(), mock.AnythingOfType("string")).
		Once().
		Return([]PlayedGame{}, nil)

	players, err := service.GetAllPlayedGames(t.Context(), "player-id")
	assert.NoError(t, err)
	assert.NotEqual(t, nil, players)
}

func TestGetOnePlayedGame(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	service := NewService(
		playerRepository,
		playedGameRepository,
		NewMockGameService(t),
	)

	playedGameRepository.
		On("FindOne", t.Context(), mock.AnythingOfType("int")).
		Once().
		Return(&PlayedGame{}, nil)

	_, err := service.GetOnePlayedGame(t.Context(), 1)
	assert.NoError(t, err)
}

func TestCreatePlayedGame(t *testing.T) {
	gameID := 2
	newID := 1
	played := []PlayedGame{
		validPlayedGame(),
	}
	played[0].Status = PlayedGameStatusCompleted

	playerRepository := NewMockPlayerRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)
	gameService := NewMockGameService(t)

	service := NewService(
		playerRepository,
		playedGameRepository,
		gameService,
	)

	gameService.
		On("GetOne", t.Context(), mock.Anything).
		Once().
		Return(&game.Game{}, nil)

	playedGameRepository.
		On("FindAll", t.Context(), mock.AnythingOfType("string")).
		Once().
		Return(played, nil)

	playedGameRepository.
		On("Insert", t.Context(), mock.AnythingOfType("*player.PlayedGame")).
		Once().
		Return(newID, nil)

	playerID, err := service.CreatePlayedGame(t.Context(), uuid.NewString(), gameID)
	assert.NoError(t, err)
	assert.Equal(t, newID, playerID)
}

func TestUpdatePlayedGame(t *testing.T) {
	validPlayer := validPlayer()
	validPlayedGame := validPlayedGame()
	validPlayedGame.PlayerID = validPlayer.ID

	playerRepository := NewMockPlayerRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)
	gameService := NewMockGameService(t)

	service := NewService(
		playerRepository,
		playedGameRepository,
		gameService,
	)

	playedGameRepository.
		On("FindOne", t.Context(), mock.AnythingOfType("int")).
		Once().
		Return(&validPlayedGame, nil)

	var upd PlayedGame
	playedGameRepository.
		On("Update", t.Context(), mock.MatchedBy(func(pg *PlayedGame) bool {
			if pg.PlayerID != validPlayer.ID {
				return false
			}
			upd = *pg
			return true
		})).
		Once().
		Return(validPlayedGame.ID, nil)

	newStatus := PlayedGameStatusCompleted
	newStartedAt := time.Now().Add(-2 * time.Hour)
	newCompletedAt := time.Now()
	newPlayTime := types.NewDurationString(time.Hour)
	newRating := 100
	newComment := testutil.Faker().Blurb()

	params := PlayedGameUpdate{
		Comment:     &newComment,
		Rating:      &newRating,
		Status:      &newStatus,
		StartedAt:   &newStartedAt,
		CompletedAt: &newCompletedAt,
		PlayTime:    &newPlayTime,
	}
	playedID, err := service.UpdatePlayedGame(
		t.Context(),
		validPlayer.ID,
		validPlayedGame.ID,
		params,
	)
	assert.NoError(t, err)

	assert.Equal(t, playedID, validPlayedGame.ID)
	assert.Equal(t, newComment, *upd.Comment)
	assert.Equal(t, newRating, *upd.Rating)
	assert.Equal(t, newStatus, upd.Status)
	assert.Equal(t, newStartedAt, upd.StartedAt)
	assert.Equal(t, newCompletedAt, *upd.CompletedAt)
	assert.Equal(t, newPlayTime, *upd.PlayTime)

	// should not change
	assert.Equal(t, validPlayedGame.Points, upd.Points)
}

func TestUpdatePlayedGame_ConsecutiveDrop(t *testing.T) {
	validPlayer := validPlayer()
	previousPlayedGame := validPlayedGame()
	validPlayedGame := validPlayedGame()

	previousPlayedGame.PlayerID = validPlayer.ID
	previousPlayedGame.Status = PlayedGameStatusDropped
	previousPlayedGame.Points = -1
	validPlayedGame.PlayerID = validPlayer.ID
	validPlayedGame.StartedAt = time.Now().Add(-3 * time.Hour)

	playerRepository := NewMockPlayerRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)
	gameService := NewMockGameService(t)

	service := NewService(
		playerRepository,
		playedGameRepository,
		gameService,
	)

	playedGameRepository.
		On("FindOne", t.Context(), mock.AnythingOfType("int")).
		Once().
		Return(&validPlayedGame, nil)

	playedGameRepository.
		On("FindLastNotReroll", t.Context(), validPlayer.ID).
		Once().
		Return(&previousPlayedGame, nil)

	var upd PlayedGame
	playedGameRepository.
		On("Update", t.Context(), mock.MatchedBy(func(pg *PlayedGame) bool {
			if pg.PlayerID != validPlayer.ID {
				return false
			}
			upd = *pg
			return true
		})).
		Once().
		Return(validPlayedGame.ID, nil)

	newStatus := PlayedGameStatusDropped

	params := PlayedGameUpdate{
		Status: &newStatus,
	}
	playedID, err := service.UpdatePlayedGame(
		t.Context(),
		validPlayer.ID,
		validPlayedGame.ID,
		params,
	)
	assert.NoError(t, err)

	assert.Equal(t, playedID, validPlayedGame.ID)
	assert.Equal(t, -2, upd.Points)
}
