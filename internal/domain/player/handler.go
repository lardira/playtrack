package player

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/domain/game"
	"github.com/lardira/playtrack/internal/pkg/password"
)

type PlayerRepository interface {
	FindAll(context.Context) ([]Player, error)
	FindOne(ctx context.Context, id string) (*Player, error)
	Insert(context.Context, *Player) (string, error)
	FindAllPlayedGames(ctx context.Context, playerID string) ([]PlayedGame, error)
	FindOnePlayedGame(ctx context.Context, playerID string, id int) (*PlayedGame, error)
	InsertPlayedGame(ctx context.Context, player *PlayedGame) (int, error)
}

type GameRepository interface {
	FindOne(ctx context.Context, id int) (*game.Game, error)
}

type Handler struct {
	playerRepository PlayerRepository
	gameRepository   GameRepository
}

func NewHandler(playerRepository PlayerRepository, gameRepository GameRepository) *Handler {
	return &Handler{
		playerRepository: playerRepository,
		gameRepository:   gameRepository,
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
	huma.Patch(grp, "/{id}", domain.EndpointNotImplemented)
	huma.Get(grp, "/{id}/played-games", h.GetAllPlayedGames)
	huma.Get(grp, "/{id}/played-games/{gameID}", h.GetOnePlayedGame)
	huma.Post(grp, "/{id}/played-games", h.CreatePlayedGame)
	huma.Patch(grp, "/{id}/played-games", domain.EndpointNotImplemented)
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

func (h *Handler) GetAllPlayedGames(ctx context.Context, i *struct {
	ID string `path:"id" format:"uuid"`
}) (*domain.ResponseItems[PlayedGame], error) {
	games, err := h.playerRepository.FindAllPlayedGames(ctx, i.ID)
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
	game, err := h.playerRepository.FindOnePlayedGame(ctx, i.ID, i.GameID)
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

	id, err := h.playerRepository.InsertPlayedGame(ctx, &nPlayed)
	if err != nil {
		return nil, huma.Error500InternalServerError("create", err)
	}

	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) containsNonterminatedPlayed(ctx context.Context, playerID string) error {
	allPlayed, err := h.playerRepository.FindAllPlayedGames(ctx, playerID)
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
