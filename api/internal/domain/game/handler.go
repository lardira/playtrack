package game

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
)

type Handler struct {
	gameService *Service
}

func NewHandler(gameService *Service) *Handler {
	return &Handler{
		gameService: gameService,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/games")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"games"}
	})

	huma.Register(grp, huma.Operation{
		OperationID: "games-get-all",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "get all games",
		Description: "get all games",
	}, h.GetAll)

	huma.Register(grp, huma.Operation{
		OperationID: "games-get-one",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "get game",
		Description: "get one game",
	}, h.GetOne)

	huma.Register(grp, huma.Operation{
		OperationID: "games-post-create",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "create game",
		Description: "create a new game",
	}, h.Create)
}

func (h *Handler) GetAll(ctx context.Context, i *struct{}) (*domain.ResponseItems[Game], error) {
	games, err := h.gameService.GetAll(ctx)
	if err != nil {
		log.Printf("game find all: %v", err)
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[Game]{}
	resp.Body.Items = games
	return &resp, nil
}

func (h *Handler) GetOne(ctx context.Context, i *struct {
	ID int `path:"id"`
}) (*domain.ResponseItem[Game], error) {
	game, err := h.gameService.GetOne(ctx, i.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("find", err)
	}

	resp := domain.ResponseItem[Game]{}
	resp.Body.Item = game
	return &resp, nil
}

func (h *Handler) Create(
	ctx context.Context,
	i *RequestCreateGame,
) (*domain.ResponseID[int], error) {
	id, err := h.gameService.CreateGame(ctx, i)
	if err != nil {
		return nil, huma.Error500InternalServerError("create", err)
	}

	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}
