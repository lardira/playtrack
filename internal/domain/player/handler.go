package player

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/pkg/password"
)

type PlayerRepository interface {
	FindAll(context.Context) ([]Player, error)
	FindOne(ctx context.Context, id string) (*Player, error)
	Insert(context.Context, *Player) (string, error)
}

type Handler struct {
	playerRepository PlayerRepository
}

func NewHandler(playerRepository PlayerRepository) *Handler {
	return &Handler{
		playerRepository: playerRepository,
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
	huma.Patch(grp, "/{id}", h.GetOne)
	huma.Get(grp, "/{id}/games-played", h.GetOne)
	huma.Post(grp, "/{id}/games-played", h.GetOne)
	huma.Patch(grp, "/{id}/games-played", h.GetOne)
}

func (h *Handler) GetAll(ctx context.Context, i *struct{}) (*domain.ResponseItems[Player], error) {
	games, err := h.playerRepository.FindAll(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[Player]{}
	resp.Body.Items = games
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
