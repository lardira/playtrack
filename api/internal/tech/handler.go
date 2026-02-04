package tech

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	checker *HealthChecker
}

func NewHandler(checker *HealthChecker) *Handler {
	return &Handler{
		checker: checker,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/tech")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"tech"}
	})

	huma.Get(grp, "/health", func(ctx context.Context, i *struct{}) (*HealthResponse, error) {
		resp := HealthResponse{}
		resp.Body.Status = Status{
			DB:     h.checker.Ok(),
			Server: true,
		}

		return &resp, nil
	})
}
