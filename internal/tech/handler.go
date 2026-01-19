package tech

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/tech")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"tech"}
	})

	huma.Get(grp, "/health", func(ctx context.Context, i *struct{}) (*HealthResponse, error) {
		// TODO: check db and s3
		resp := HealthResponse{}
		resp.Body.Message = "ok"
		return &resp, nil
	})
}
