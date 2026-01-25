package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/domain/game"
	"github.com/lardira/playtrack/internal/pkg/password"
)

type PlayerRepository interface {
	FindAll(context.Context) ([]Player, error)
	FindOne(ctx context.Context, id string) (*Player, error)
	Insert(context.Context, *Player) (string, error)
	Update(ctx context.Context, player *PlayerUpdate) (string, error)
}

type PlayedGameRepository interface {
	FindAll(ctx context.Context, playerID string) ([]PlayedGame, error)
	FindOne(ctx context.Context, playerID string, id int) (*PlayedGame, error)
	FindLast(ctx context.Context, playerID string) (*PlayedGame, error)
	Insert(ctx context.Context, player *PlayedGame) (int, error)
	Update(ctx context.Context, game *PlayedGameUpdate) (string, error)
}

type GameRepository interface {
	FindOne(ctx context.Context, id int) (*game.Game, error)
}

type Handler struct {
	playerRepository     PlayerRepository
	playedGameRepository PlayedGameRepository
	gameRepository       GameRepository
}

func NewHandler(
	playerRepository PlayerRepository,
	gameRepository GameRepository,
	playedGameRepository PlayedGameRepository,
) *Handler {
	return &Handler{
		playerRepository:     playerRepository,
		playedGameRepository: playedGameRepository,
		gameRepository:       gameRepository,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/players")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"players"}
	})

	huma.Get(grp, "/", h.GetAll)
	// huma.Get(grp, "/leaderboard", h.GetAll)
	huma.Get(grp, "/{id}", h.GetOne)
	huma.Post(grp, "/", h.Create)
	huma.Patch(grp, "/{id}", h.Update)
	huma.Get(grp, "/{id}/played-games", h.GetAllPlayedGames)
	huma.Get(grp, "/{id}/played-games/{gameID}", h.GetOnePlayedGame)
	huma.Post(grp, "/{id}/played-games", h.CreatePlayedGame)
	huma.Patch(grp, "/{id}/played-games/{gameID}", h.UpdatePlayedGame)
}

func (h *Handler) GetAll(ctx context.Context, i *struct{}) (*domain.ResponseItems[Player], error) {
	players, err := h.playerRepository.FindAll(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[Player]{}
	resp.Body.Items = players
	return &resp, nil
}

func (h *Handler) GetOne(ctx context.Context, i *struct {
	ID string `path:"id" format:"uuid"`
}) (*domain.ResponseItem[Player], error) {
	player, err := h.playerRepository.FindOne(ctx, i.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("find", err)
	}

	resp := domain.ResponseItem[Player]{}
	resp.Body.Item = player
	return &resp, nil
}

func (h *Handler) Create(
	ctx context.Context,
	i *RequestCreatePlayer,
) (*domain.ResponseID[string], error) {
	nPlayer := Player{
		Username: i.Body.Username,
		Img:      i.Body.Img,
		Email:    i.Body.Email,
		Password: i.Body.Password,
	}
	if err := nPlayer.Valid(); err != nil {
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	hashedPassword, err := password.Hash(nPlayer.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("could not create player")
	}
	nPlayer.Password = hashedPassword

	id, err := h.playerRepository.Insert(ctx, &nPlayer)
	if err != nil {
		return nil, huma.Error500InternalServerError("create", err)
	}

	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) Update(
	ctx context.Context,
	i *RequestUpdatePlayer,
) (*domain.ResponseID[string], error) {
	nPlayer := PlayerUpdate{
		ID:       i.PlayerID,
		Username: i.Body.Username,
		Img:      i.Body.Img,
		Email:    i.Body.Email,
	}
	if err := nPlayer.Valid(); err != nil {
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	id, err := h.playerRepository.Update(ctx, &nPlayer)
	if err != nil {
		return nil, huma.Error500InternalServerError("update", err)
	}

	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) GetAllPlayedGames(ctx context.Context, i *struct {
	ID string `path:"id" format:"uuid"`
}) (*domain.ResponseItems[PlayedGame], error) {
	games, err := h.playedGameRepository.FindAll(ctx, i.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[PlayedGame]{}
	resp.Body.Items = games
	return &resp, nil
}

func (h *Handler) GetOnePlayedGame(ctx context.Context, i *struct {
	ID     string `path:"id" format:"uuid"`
	GameID int    `path:"gameID"`
}) (*domain.ResponseItem[PlayedGame], error) {
	game, err := h.playedGameRepository.FindOne(ctx, i.ID, i.GameID)
	if err != nil {
		return nil, huma.Error500InternalServerError("find", err)
	}

	resp := domain.ResponseItem[PlayedGame]{}
	resp.Body.Item = game
	return &resp, nil
}

func (h *Handler) CreatePlayedGame(
	ctx context.Context,
	i *RequestCreatePlayedGame,
) (*domain.ResponseID[int], error) {
	game, err := h.gameRepository.FindOne(ctx, i.Body.GameID)
	if err != nil {
		return nil, huma.Error400BadRequest("game find", err)
	}

	nPlayed := PlayedGame{
		PlayerID: i.PlayerID,
		GameID:   i.Body.GameID,
		Points:   game.Points,
	}

	if err := nPlayed.Valid(); err != nil {
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	if err := h.containsNonterminatedPlayed(ctx, i.PlayerID); err != nil {
		return nil, err
	}

	id, err := h.playedGameRepository.Insert(ctx, &nPlayed)
	if err != nil {
		return nil, huma.Error500InternalServerError("create", err)
	}

	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) UpdatePlayedGame(
	ctx context.Context,
	i *RequestUpdatePlayedGame,
) (*domain.ResponseID[string], error) {
	nGame := PlayedGameUpdate{
		ID:          i.GameID,
		Points:      i.Body.Points,
		Comment:     i.Body.Comment,
		Rating:      i.Body.Rating,
		Status:      i.Body.Status,
		CompletedAt: i.Body.CompletedAt,
		PlayTime:    i.Body.PlayTime,
	}
	if err := nGame.Valid(); err != nil {
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	playedGame, err := h.playedGameRepository.FindOne(ctx, i.PlayerID, i.GameID)
	if err != nil {
		return nil, huma.Error400BadRequest("entity is not found", err)
	}

	if nGame.Status != nil {
		status := *nGame.Status
		if err := playedGame.StatusNextValid(status); err != nil {
			return nil, huma.Error400BadRequest("entity is not valid", err)
		}

		newPoints := 0
		now := time.Now()
		switch status {
		case PlayedGameStatusDropped:
			newPoints = -1

			prevGame, err := h.playedGameRepository.FindLast(ctx, i.PlayerID)
			if err != nil && !errors.Is(err, ErrPlayedGameNotFound) {
				return nil, huma.Error400BadRequest("game played find", err)
			}

			// consecutive drops are stacked
			if err == nil && prevGame.Status == PlayedGameStatusDropped {
				newPoints = prevGame.Points - 1
			}

			nGame.Points = &newPoints
			nGame.CompletedAt = &now

		case PlayedGameStatusRerolled:
			nGame.Points = &newPoints
			nGame.CompletedAt = &now
		}
	}

	id, err := h.playedGameRepository.Update(ctx, &nGame)
	if err != nil {
		return nil, huma.Error500InternalServerError("update", err)
	}

	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) containsNonterminatedPlayed(ctx context.Context, playerID string) error {
	allPlayed, err := h.playedGameRepository.FindAll(ctx, playerID)
	if err != nil {
		return huma.Error400BadRequest("find played games", err)
	}

	for _, p := range allPlayed {
		if !p.StatusTerminated() {
			msg := fmt.Sprintf("player has game in nonterminated status: %v", p.ID)
			return huma.Error400BadRequest(msg)
		}
	}
	return nil
}
