package game

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/games")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"games"}
	})

	huma.Get(grp, "/", domain.EndpointNotImplemented)
	huma.Get(grp, "/{id}", domain.EndpointNotImplemented)
	huma.Post(grp, "/", domain.EndpointNotImplemented)
}
