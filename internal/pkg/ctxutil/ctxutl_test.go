package ctxutil

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
)

func TestGetSetPlayerID(t *testing.T) {
	playerID := uuid.NewString()
	ctx := context.Background()

	ctx = SetPlayerID(ctx, playerID)
	assert.NotEqual(t, nil, ctx)

	parsed, ok := GetPlayerID(ctx)
	assert.True(t, ok)
	assert.Equal(t, playerID, parsed)
}

func TestGetSetPlayerID_IDNotSet(t *testing.T) {
	ctx := context.Background()

	parsed, ok := GetPlayerID(ctx)
	assert.False(t, ok)
	assert.Zero(t, parsed)
}
