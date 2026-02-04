package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/pkg/apiutil"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
	"github.com/lardira/playtrack/internal/pkg/testutil"
)

const (
	testSecret = "test"
)

type TestHumaContext interface {
	huma.Context
}

type testCtx struct {
	TestHumaContext

	onSetStatus func(int)
	onHeader    func() string
}

func (t testCtx) Header(_ string) string {
	return t.onHeader()
}

func (t testCtx) SetStatus(code int) {
	t.onSetStatus(code)
}

func (t testCtx) Context() context.Context {
	return context.Background()
}

func TestAuthMiddleware(t *testing.T) {
	playerID := uuid.NewString()
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   playerID,
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		Audience:  []string{apiutil.RolePlayer},
	}
	token := jwt.NewWithClaims(apiutil.DefaultSigningMethod, claims)
	signedToken, _ := token.SignedString([]byte(testSecret))

	authFunc := Authorize(testSecret)

	ctx := testCtx{
		onHeader: func() string {
			return authPrefix + signedToken
		},
		onSetStatus: func(code int) {
			assert.NotEqual(t, http.StatusUnauthorized, code)
		},
	}
	authFunc(ctx, func(ctx huma.Context) {
		authCtx, ok := ctx.(*authContext)
		assert.True(t, ok)
		assert.NotZero(t, authCtx)

		ctxP, ok := ctxutil.GetPlayer(authCtx.Context())
		assert.True(t, ok)
		assert.Equal(t, playerID, ctxP.ID)
	})
}

func TestAuthMiddleware_Admin(t *testing.T) {
	playerID := uuid.NewString()
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   playerID,
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		Audience:  []string{apiutil.RolePlayer, apiutil.RoleAdmin},
	}
	token := jwt.NewWithClaims(apiutil.DefaultSigningMethod, claims)
	signedToken, _ := token.SignedString([]byte(testSecret))

	authFunc := Authorize(testSecret)

	ctx := testCtx{
		onHeader: func() string {
			return authPrefix + signedToken
		},
		onSetStatus: func(code int) {
			assert.NotEqual(t, http.StatusUnauthorized, code)
		},
	}
	authFunc(ctx, func(ctx huma.Context) {
		authCtx, ok := ctx.(*authContext)
		assert.True(t, ok)
		assert.NotZero(t, authCtx)

		ctxP, ok := ctxutil.GetPlayer(authCtx.Context())
		assert.True(t, ok)
		assert.Equal(t, playerID, ctxP.ID)
		assert.True(t, ctxP.IsAdmin)
	})
}

func TestAuthMiddleware_InvalidHeader(t *testing.T) {
	tcases := []struct {
		name     string
		onHeader func() string
		code     int
	}{
		{
			"empty header",
			func() string { return "" },
			http.StatusUnauthorized,
		},
		{
			"empty bearer",
			func() string { return authPrefix },
			http.StatusUnauthorized,
		},
		{
			"bearer invalid token",
			func() string { return testutil.Faker().Regex("a{10}\\.b{10}\\.c{10}") },
			http.StatusUnauthorized,
		},
	}

	authFunc := Authorize(testSecret)

	for _, tt := range tcases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testCtx{
				onSetStatus: func(code int) {
					assert.Equal(t, tt.code, code)
				},
				onHeader: tt.onHeader,
			}

			authFunc(ctx, func(ctx huma.Context) {
				authCtx, ok := ctx.(*authContext)
				assert.True(t, ok)
				assert.NotZero(t, authCtx)

				_, ok = ctxutil.GetPlayer(authCtx.Context())
				assert.False(t, ok)
			})
		})
	}
}
