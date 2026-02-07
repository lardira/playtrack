package game

import (
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/lardira/playtrack/internal/pkg/testutil"
	mock "github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	gameRepository := NewMockGameRepository(t)
	gameService := NewService(gameRepository)
	handler := NewHandler(gameService)

	games := make([]Game, 2)
	testutil.Faker().Struct(&games[0])
	testutil.Faker().Struct(&games[1])

	gameRepository.
		On("FindAll", t.Context()).
		Once().
		Return(games, nil)

	resp, err := handler.GetAll(t.Context(), nil)
	assert.NoError(t, err)
	assert.Equal(t, games, resp.Body.Items)
}

func TestGetOne(t *testing.T) {
	gameRepository := NewMockGameRepository(t)
	gameService := NewService(gameRepository)
	handler := NewHandler(gameService)

	var game Game
	testutil.Faker().Struct(&game)

	gameRepository.
		On("FindOne", t.Context(), game.ID).
		Once().
		Return(&game, nil)

	req := struct {
		ID int `path:"id"`
	}{
		ID: game.ID,
	}

	resp, err := handler.GetOne(t.Context(), &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp.Body.Item)
	assert.Equal(t, game, *resp.Body.Item)
}

func TestGetOne_NotFound(t *testing.T) {
	gameRepository := NewMockGameRepository(t)
	gameService := NewService(gameRepository)
	handler := NewHandler(gameService)

	var game Game
	testutil.Faker().Struct(&game)

	gameRepository.
		On("FindOne", t.Context(), game.ID).
		Once().
		Return(nil, errors.New("not found"))

	req := struct {
		ID int `path:"id"`
	}{
		ID: game.ID,
	}

	resp, err := handler.GetOne(t.Context(), &req)
	assert.Error(t, err)
	assert.Equal(t, nil, resp)
}

func TestGetCreate(t *testing.T) {
	gameRepository := NewMockGameRepository(t)
	gameService := NewService(gameRepository)
	handler := NewHandler(gameService)

	newID := testutil.Faker().Int()
	hoursToBeat := 2
	url := testutil.Faker().URL()

	gameRepository.
		On("Insert", t.Context(), mock.AnythingOfType("*game.Game")).
		Once().
		Return(newID, nil)

	var req RequestCreateGame
	req.Body.HoursToBeat = hoursToBeat
	req.Body.Title = testutil.Faker().MovieName()
	req.Body.URL = &url

	resp, err := handler.Create(t.Context(), &req)
	assert.NoError(t, err)
	assert.Equal(t, newID, resp.Body.ID)
}
