package game

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
)

type GameRepository interface {
	FindAll(context.Context) ([]Game, error)
	FindOne(ctx context.Context, id int) (*Game, error)
	Insert(context.Context, *Game) (int, error)
}

type Handler struct {
	gameRepository GameRepository
}

func NewHandler(gameRepository GameRepository) *Handler {
	return &Handler{
		gameRepository: gameRepository,
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
	games, err := h.gameRepository.FindAll(ctx)
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
	game, err := h.gameRepository.FindOne(ctx, i.ID)
	if err != nil {
		log.Printf("game find one: %v", err)
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
	nGame := Game{
		HoursToBeat: i.Body.HoursToBeat,
		Title:       i.Body.Title,
		URL:         i.Body.URL,
	}
	nGame.CalculatePoints()

	if err := nGame.Valid(); err != nil {
		log.Printf("game valid: %v", err)
		return nil, huma.Error400BadRequest("game is not valid", err)
	}

	id, err := h.gameRepository.Insert(ctx, &nGame)
	if err != nil {
		log.Printf("game insert: %v", err)
		return nil, huma.Error500InternalServerError("create", err)
	}

	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}
