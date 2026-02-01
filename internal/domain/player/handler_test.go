package player

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/domain/game"
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
		Return(player.ID, nil)

	req := RequestUpdatePlayer{}
	req.PlayerID = player.ID
	req.Body.Username = &player.Username
	req.Body.Img = player.Img
	req.Body.Email = player.Email

	resp, err := handler.Update(ctx, &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, player.ID, resp.Body.ID)
}

func TestUpdate_OtherPlayer(t *testing.T) {
	player := validPlayer()
	otherPlayer := validPlayer()
	ctx := ctxutil.SetPlayerID(t.Context(), otherPlayer.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playerRepository.AssertNotCalled(t, "Update")

	req := RequestUpdatePlayer{}
	req.PlayerID = player.ID
	req.Body.Username = &player.Username
	req.Body.Img = player.Img
	req.Body.Email = player.Email

	_, err := handler.Update(ctx, &req)
	assert.Error(t, err)
}

func TestGetAllPlayedGames(t *testing.T) {
	playerID := uuid.NewString()
	games := make([]PlayedGame, 2)
	testutil.Faker().Struct(&games[0])
	testutil.Faker().Struct(&games[1])

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playedGameRepository.
		On("FindAll", t.Context(), playerID).
		Once().
		Return(games, nil)

	req := struct {
		ID string `path:"id" format:"uuid"`
	}{
		ID: playerID,
	}

	resp, err := handler.GetAllPlayedGames(t.Context(), &req)
	assert.NoError(t, err)
	assert.Equal(t, games, resp.Body.Items)
}

func TestGetOnePlayedGame(t *testing.T) {
	game := validPlayedGame()

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playedGameRepository.
		On("FindOne", t.Context(), game.PlayerID, game.ID).
		Once().
		Return(&game, nil)

	req := struct {
		PlayerID string `path:"id" format:"uuid"`
		GameID   int    `path:"gameID"`
	}{
		PlayerID: game.PlayerID,
		GameID:   game.ID,
	}

	resp, err := handler.GetOnePlayedGame(t.Context(), &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, game, *resp.Body.Item)
}

func TestCreatePlayedGame(t *testing.T) {
	game := game.Game{
		ID:          testutil.Faker().Int(),
		Points:      2,
		HoursToBeat: 3,
		Title:       testutil.Faker().MovieName(),
		CreatedAt:   time.Now(),
	}
	player := validPlayer()
	played := validPlayedGame()
	played.PlayerID = player.ID
	playedGame := []PlayedGame{played}
	ctx := ctxutil.SetPlayerID(t.Context(), player.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	gameRepository.
		On("FindOne", ctx, game.ID).
		Once().
		Return(&game, nil)

	playedGameRepository.
		On("FindAll", ctx, player.ID).
		Once().
		Return(playedGame, nil)

	playedGameRepository.
		On("Insert", ctx, mock.MatchedBy(func(p *PlayedGame) bool {
			if p.Points != game.Points {
				return false
			}
			if p.GameID != game.ID {
				return false
			}
			if p.PlayerID != player.ID {
				return false
			}
			return true
		})).
		Once().
		Return(played.ID, nil)

	req := RequestCreatePlayedGame{}
	req.PlayerID = player.ID
	req.Body.GameID = game.ID

	resp, err := handler.CreatePlayedGame(ctx, &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, played.ID, resp.Body.ID)
}

func TestCreatePlayed_OtherPlayer(t *testing.T) {
	player := validPlayer()
	otherPlayer := validPlayer()

	ctx := ctxutil.SetPlayerID(t.Context(), otherPlayer.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	gameRepository.AssertNotCalled(t, "FindOne")
	playedGameRepository.AssertNotCalled(t, "FindAll")
	playedGameRepository.AssertNotCalled(t, "Insert")

	req := RequestCreatePlayedGame{}
	req.PlayerID = player.ID
	req.Body.GameID = testutil.Faker().Int()

	resp, err := handler.CreatePlayedGame(ctx, &req)
	assert.Error(t, err)
	assert.Equal(t, nil, resp)
}

func TestUpdatePlayedGame(t *testing.T) {
	player := validPlayer()
	played := validPlayedGame()
	played.PlayerID = player.ID
	ctx := ctxutil.SetPlayerID(t.Context(), player.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playedGameRepository.
		On("FindOne", ctx, player.ID, played.ID).
		Once().
		Return(&played, nil)

	playedGameRepository.
		On("Update", ctx, mock.AnythingOfType("*player.PlayedGameUpdate")).
		Once().
		Return(played.ID, nil)

	req := RequestUpdatePlayedGame{}
	req.PlayerID = player.ID
	req.GameID = played.ID
	req.Body.Rating = played.Rating

	resp, err := handler.UpdatePlayedGame(ctx, &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, played.ID, resp.Body.ID)
}

func TestUpdatePlayedGame_OtherPlayer(t *testing.T) {
	player := validPlayer()
	otherPlayer := validPlayer()
	played := validPlayedGame()
	ctx := ctxutil.SetPlayerID(t.Context(), otherPlayer.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playedGameRepository.AssertNotCalled(t, "FindOne")
	playedGameRepository.AssertNotCalled(t, "Update")

	req := RequestUpdatePlayedGame{}
	req.PlayerID = player.ID
	req.GameID = played.ID
	req.Body.Rating = played.Rating

	resp, err := handler.UpdatePlayedGame(ctx, &req)
	assert.Error(t, err)
	assert.Equal(t, nil, resp)
}

func TestUpdatePlayedGame_ConsecutiveDrop(t *testing.T) {
	played := []PlayedGame{
		validPlayedGame(),
		validPlayedGame(),
	}

	player := validPlayer()
	played[0].PlayerID = player.ID
	played[0].Status = PlayedGameStatusDropped
	played[0].Points = -1
	played[1].PlayerID = player.ID
	played[1].Status = PlayedGameStatusInProgress
	played[1].Points = -2

	ctx := ctxutil.SetPlayerID(t.Context(), player.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playedGameRepository.
		On("FindOne", ctx, player.ID, played[1].ID).
		Once().
		Return(&played[1], nil)

	playedGameRepository.
		On("FindLastNotReroll", ctx, player.ID).
		Once().
		Return(&played[0], nil)

	playedGameRepository.
		On("Update", ctx, mock.MatchedBy(func(p *PlayedGameUpdate) bool {
			if p.Points == nil || *p.Points != played[1].Points {
				return false
			}
			if p.CompletedAt == nil {
				return false
			}
			return true
		})).
		Once().
		Return(played[1].ID, nil)

	newStatus := PlayedGameStatusDropped
	req := RequestUpdatePlayedGame{}
	req.PlayerID = player.ID
	req.GameID = played[1].ID
	req.Body.Rating = played[1].Rating
	req.Body.Status = &newStatus

	resp, err := handler.UpdatePlayedGame(ctx, &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, played[1].ID, resp.Body.ID)
}

func TestUpdatePlayedGame_Reroll(t *testing.T) {
	played := []PlayedGame{
		validPlayedGame(),
		validPlayedGame(),
	}

	player := validPlayer()
	played[0].PlayerID = player.ID
	played[0].Status = PlayedGameStatusRerolled
	played[0].Points = 0
	played[1].PlayerID = player.ID
	played[1].Status = PlayedGameStatusInProgress
	played[1].Points = 0

	ctx := ctxutil.SetPlayerID(t.Context(), player.ID)

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	playedGameRepository.
		On("FindOne", ctx, player.ID, played[1].ID).
		Once().
		Return(&played[1], nil)

	playedGameRepository.
		On("Update", ctx, mock.MatchedBy(func(p *PlayedGameUpdate) bool {
			if p.Points == nil || *p.Points != played[1].Points {
				return false
			}
			if p.CompletedAt == nil {
				return false
			}
			return true
		})).
		Once().
		Return(played[1].ID, nil)

	newStatus := PlayedGameStatusRerolled
	req := RequestUpdatePlayedGame{}
	req.PlayerID = player.ID
	req.GameID = played[1].ID
	req.Body.Rating = played[1].Rating
	req.Body.Status = &newStatus

	resp, err := handler.UpdatePlayedGame(ctx, &req)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, resp)
	assert.Equal(t, played[1].ID, resp.Body.ID)
}

func TestContainsNonterminatedPlayed(t *testing.T) {
	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	played := []PlayedGame{
		validPlayedGame(),
		validPlayedGame(),
	}

	player := validPlayer()
	played[0].PlayerID = player.ID
	played[1].PlayerID = player.ID

	playedGameRepository.
		On("FindAll", t.Context(), player.ID).
		Once().
		Return(played, nil)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	err := handler.containsNonterminatedPlayed(t.Context(), player.ID)
	assert.NoError(t, err)
}

func TestContainsNonterminatedPlayed_Nonterminated(t *testing.T) {

	playerRepository := NewMockPlayerRepository(t)
	gameRepository := NewMockGameRepository(t)
	playedGameRepository := NewMockPlayedGameRepository(t)

	played := []PlayedGame{
		validPlayedGame(),
		validPlayedGame(),
	}

	player := validPlayer()
	played[0].PlayerID = player.ID
	played[1].PlayerID = player.ID
	played[1].Status = PlayedGameStatusAdded

	playedGameRepository.
		On("FindAll", t.Context(), player.ID).
		Once().
		Return(played, nil)

	handler := NewHandler(playerRepository, gameRepository, playedGameRepository)

	err := handler.containsNonterminatedPlayed(t.Context(), player.ID)
	assert.Error(t, err)
}
