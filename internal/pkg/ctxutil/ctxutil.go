package ctxutil

import (
	"context"
)

type contextKey string

const (
	keyPlayerID contextKey = "playerID"
)

func GetPlayerID(ctx context.Context) (string, bool) {
	v := ctx.Value(keyPlayerID)
	playerID, ok := v.(string)
	return playerID, ok
}

func SetPlayerID(ctx context.Context, playerID string) context.Context {
	return context.WithValue(ctx, keyPlayerID, playerID)
}
