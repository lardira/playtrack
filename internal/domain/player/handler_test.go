package player

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
	"github.com/lardira/playtrack/internal/pkg/testutil"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	players := make([]Player, 2)
	testutil.Faker().Struct(&players[0])
	testutil.Faker().Struct(&players[1])

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playerRepository.
		On("FindAll", t.Context()).
		Once().
		Return(players, nil)

	resp, err := handler.GetAll(t.Context(), nil)
	assert.NoError(t, err)
	assert.Equal(t, players, resp.Body.Items)
}

func TestGetOne(t *testing.T) {
	player := validPlayer()

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playerRepository.
		On("FindOne", t.Context(), mock.AnythingOfType("string")).
		Once().
		Return(&player, nil)

	req := struct {
		ID string `path:"id" format:"uuid"`
	}{
		ID: player.ID,
	}

	resp, err := handler.GetOne(t.Context(), &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, player, *resp.Body.Item)
}

func TestUpdate(t *testing.T) {
	player := validPlayer()
	ctx := ctxutil.SetPlayerID(t.Context(), player.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playerRepository.
		On("Update", ctx, mock.AnythingOfType("*player.PlayerUpdate")).
		Once().
		Return(&player.ID, nil)

	req := struct {
		ID string `path:"id" format:"uuid"`
	}{
		ID: player.ID,
	}

	resp, err := handler.GetOne(ctx, &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, player, *resp.Body.Item)
}
