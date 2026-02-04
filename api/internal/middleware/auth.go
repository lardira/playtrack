package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/pkg/apiutil"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
)

const (
	authHeader = "Authorization"

	authPrefix = "Bearer "
)

type humaContext huma.Context

type authContext struct {
	humaContext
	playerID string
	isAdmin  bool
}

func (c *authContext) Context() context.Context {
	return ctxutil.SetPlayer(
		c.humaContext.Context(),
		ctxutil.CtxPlayer{ID: c.playerID, IsAdmin: c.isAdmin},
	)
}

func Authorize(secret string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		auth := ctx.Header(authHeader)
		tokenString, ok := strings.CutPrefix(auth, authPrefix)
		if !ok {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		parseToken := func(t *jwt.Token) (any, error) {
			exp, err := t.Claims.GetExpirationTime()
			if err != nil {
				return nil, err
			}
			if exp.Before(time.Now()) {
				return nil, fmt.Errorf("token expired")
			}
			return []byte(secret), nil
		}

		token, err := jwt.Parse(
			tokenString,
			parseToken,
			jwt.WithValidMethods([]string{apiutil.DefaultSigningMethod.Alg()}),
			jwt.WithAudience(apiutil.RoleAdmin, apiutil.RolePlayer),
			jwt.WithIssuedAt(),
			jwt.WithNotBeforeRequired(),
		)
		if err != nil {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		playerID, err := token.Claims.GetSubject()
		if err != nil {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		if err := uuid.Validate(playerID); err != nil {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		aud, err := token.Claims.GetAudience()
		if err != nil {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		authCtx := authContext{
			humaContext: ctx,
			playerID:    playerID,
			isAdmin:     slices.Contains(aud, apiutil.RoleAdmin),
		}
		next(&authCtx)
	}
}
