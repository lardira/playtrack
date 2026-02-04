package ctxutil

import (
	"context"
)

type contextKey string

const (
	keyPlayer contextKey = "player"
)

type CtxPlayer struct {
	ID      string
	IsAdmin bool
}

func GetPlayer(ctx context.Context) (CtxPlayer, bool) {
	v := ctx.Value(keyPlayer)
	player, ok := v.(CtxPlayer)
	return player, ok && player.ID != ""
}

func SetPlayer(ctx context.Context, p CtxPlayer) context.Context {
	return context.WithValue(ctx, keyPlayer, p)
}
