package player

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
)

type Handler struct {
	playerService *Service
}

func NewHandler(
	playerService *Service,
) *Handler {
	return &Handler{
		playerService: playerService,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/players")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"players"}
	})

	huma.Register(grp, huma.Operation{
		OperationID: "players-get-all",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "get all players",
		Description: "get all players",
	}, h.GetAll)

	huma.Register(grp, huma.Operation{
		OperationID: "players-get-one",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "get one player",
		Description: "get one player",
	}, h.GetOne)

	huma.Register(grp, huma.Operation{
		OperationID: "players-update-one",
		Method:      http.MethodPatch,
		Path:        "/{id}",
		Summary:     "update player",
		Description: "update a player",
	}, h.Update)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-get-all",
		Method:      http.MethodGet,
		Path:        "/{id}/played-games",
		Summary:     "get all played games",
		Description: "get all played games",
	}, h.GetAllPlayedGames)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-get-one",
		Method:      http.MethodGet,
		Path:        "/{id}/played-games/{gameID}",
		Summary:     "get one played game",
		Description: "get one played game",
	}, h.GetOnePlayedGame)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-create-one",
		Method:      http.MethodPost,
		Path:        "/{id}/played-games",
		Summary:     "create played game",
		Description: "create a new played game",
	}, h.CreatePlayedGame)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-update-one",
		Method:      http.MethodPatch,
		Path:        "/{id}/played-games/{gameID}",
		Summary:     "update played game",
		Description: "update a played game",
	}, h.UpdatePlayedGame)
}

func (h *Handler) GetAll(ctx context.Context, i *struct{}) (*domain.ResponseItems[Player], error) {
	players, err := h.playerService.GetAll(ctx)
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
	player, err := h.playerService.GetOne(ctx, i.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("find", err)
	}

	resp := domain.ResponseItem[Player]{}
	resp.Body.Item = player
	return &resp, nil
}

func (h *Handler) Update(
	ctx context.Context,
	i *RequestUpdatePlayer,
) (*domain.ResponseID[string], error) {
	if ok := checkAuthorizedFor(ctx, i.PlayerID); !ok {
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	id, err := h.playerService.Update(ctx, i.PlayerID, i.Body)
	if err != nil {
		log.Printf("player update: %v", err)
		return nil, huma.Error500InternalServerError("update", err)
	}

	log.Printf("player %v updated", id)
	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) GetAllPlayedGames(ctx context.Context, i *struct {
	ID string `path:"id" format:"uuid"`
}) (*domain.ResponseItems[PlayedGame], error) {
	games, err := h.playerService.GetAllPlayedGames(ctx, i.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[PlayedGame]{}
	resp.Body.Items = games
	return &resp, nil
}

func (h *Handler) GetOnePlayedGame(ctx context.Context, i *struct {
	PlayerID string `path:"id" format:"uuid"`
	GameID   int    `path:"gameID"`
}) (*domain.ResponseItem[PlayedGame], error) {
	game, err := h.playerService.GetOnePlayedGame(ctx, i.GameID)
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
	if ok := checkAuthorizedFor(ctx, i.PlayerID); !ok {
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	id, err := h.playerService.CreatePlayedGame(ctx, i.PlayerID, i.Body.GameID)
	if err != nil {
		return nil, huma.Error400BadRequest("create", err)
	}

	log.Printf("played game %v created", id)
	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) UpdatePlayedGame(
	ctx context.Context,
	i *RequestUpdatePlayedGame,
) (*domain.ResponseID[int], error) {
	if ok := checkAuthorizedFor(ctx, i.PlayerID); !ok {
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	id, err := h.playerService.UpdatePlayedGame(ctx, i.PlayerID, i.GameID, i.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("entity is not found", err)
	}

	log.Printf("played game %v updated", id)
	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}

func checkAuthorizedFor(ctx context.Context, playerID string) bool {
	ctxPlayer, ok := ctxutil.GetPlayer(ctx)
	if !ok {
		return false
	}

	if ctxPlayer.IsAdmin {
		return true
	} else if ctxPlayer.ID != playerID {
		return false
	}
	return true
}
