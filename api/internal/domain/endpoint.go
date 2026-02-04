package domain

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func EndpointNotImplemented(ctx context.Context, i *struct{}) (*struct{}, error) {
	return nil, huma.Error501NotImplemented("not implemented yet")
}
