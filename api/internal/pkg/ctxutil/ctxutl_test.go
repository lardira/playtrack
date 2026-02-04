package ctxutil

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/pkg/testutil"
)

func TestGetSetPlayer(t *testing.T) {
	playerID := uuid.NewString()
	isAdmin := testutil.Faker().Bool()

	ctxPlayer := CtxPlayer{
		IsAdmin: isAdmin,
		ID:      playerID,
	}

	ctx := context.Background()

	ctx = SetPlayer(ctx, ctxPlayer)
	assert.NotEqual(t, nil, ctx)

	parsed, ok := GetPlayer(ctx)
	assert.True(t, ok)
	assert.Equal(t, parsed, ctxPlayer)
}

func TestGetSetPlayer_NoValue(t *testing.T) {
	ctx := context.Background()

	parsed, ok := GetPlayer(ctx)
	assert.False(t, ok)
	assert.Zero(t, parsed)
}

func TestGetSetPlayer_IDNotSet(t *testing.T) {
	ctx := SetPlayer(context.Background(), CtxPlayer{})

	parsed, ok := GetPlayer(ctx)
	assert.False(t, ok)
	assert.Zero(t, parsed)
}
