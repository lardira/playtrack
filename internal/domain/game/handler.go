package game

import (
	"context"

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

	huma.Get(grp, "/", h.GetAll)
	huma.Get(grp, "/{id}", h.GetOne)
	huma.Post(grp, "/", h.Create)
}

func (h *Handler) GetAll(ctx context.Context, i *struct{}) (*domain.ResponseItems[Game], error) {
	games, err := h.gameRepository.FindAll(ctx)
	if err != nil {
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
		Points:      i.Body.Points,
		HoursToBeat: i.Body.HoursToBeat,
		Title:       i.Body.Title,
		URL:         i.Body.URL,
	}
	if err := nGame.Valid(); err != nil {
		return nil, huma.Error400BadRequest("game is not valid", err)
	}

	id, err := h.gameRepository.Insert(ctx, &nGame)
	if err != nil {
		return nil, huma.Error500InternalServerError("create", err)
	}

	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}
