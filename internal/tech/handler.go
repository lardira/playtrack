package tech

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/tech")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"tech"}
	})

	huma.Get(grp, "/health", func(ctx context.Context, i *struct{}) (*HealthResponse, error) {
		resp := HealthResponse{}
		resp.Body.DB = true
		resp.Body.Server = true

		if err := h.pool.Ping(context.Background()); err != nil {
			resp.Body.Status.DB = false
		}

		return &resp, nil
	})
}
