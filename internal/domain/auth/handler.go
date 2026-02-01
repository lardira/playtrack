package auth

import (
	"context"
	"log"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/domain/player"
	"github.com/lardira/playtrack/internal/middleware"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
	"github.com/lardira/playtrack/internal/pkg/password"
)

const (
	defaultExpiration = 10 * 24 * time.Hour
)

type PlayerRepository interface {
	FindOneByUsername(ctx context.Context, username string) (*player.Player, error)
	Update(ctx context.Context, player *player.PlayerUpdate) (string, error)
	Insert(context.Context, *player.Player) (string, error)
}

type Handler struct {
	secret           string
	playerRepository PlayerRepository
}

func NewHandler(secret string, playerRepository PlayerRepository) *Handler {
	return &Handler{
		secret:           secret,
		playerRepository: playerRepository,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/auth")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"auth"}
	})
	secgrp := huma.NewGroup(api, "")
	secgrp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"auth"}
	})
	secgrp.UseMiddleware(
		middleware.Authorize(h.secret),
	)

	huma.Post(grp, "/register", h.RegisterPlayer)
	huma.Post(grp, "/login", h.Login)
	huma.Patch(secgrp, "/set-password", h.SetPassword)
}

func (h *Handler) Login(ctx context.Context, i *RequestLoginPlayer) (*ResponseLoginPlayer, error) {
	found, err := h.playerRepository.FindOneByUsername(ctx, i.Body.Username)
	if err != nil {
		log.Printf("login player find one: %v", err)
		return nil, huma.Error401Unauthorized("username or password is incorrect")
	}
	if !password.CompareHash(i.Body.Password, found.Password) {
		log.Printf("login compare hash: %v", err)
		return nil, huma.Error401Unauthorized("username or password is incorrect")
	}

	token, err := h.issueToken(found.ID)
	if err != nil {
		log.Printf("login issue token: %v", err)
		return nil, huma.Error500InternalServerError("could not issue token", err)
	}

	resp := ResponseLoginPlayer{}
	resp.Body.Token = token
	return &resp, nil
}

func (h *Handler) RegisterPlayer(
	ctx context.Context,
	i *RequestRegisterCreatePlayer,
) (*domain.ResponseID[string], error) {
	nPlayer := player.Player{
		Username: i.Body.Username,
		Img:      i.Body.Img,
		Email:    i.Body.Email,
		Password: i.Body.Password,
	}
	if err := nPlayer.Valid(); err != nil {
		log.Printf("register player not valid: %v", err)
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	hashedPassword, err := password.Hash(nPlayer.Password)
	if err != nil {
		log.Printf("register pass hash: %v", err)
		return nil, huma.Error500InternalServerError("could not create player")
	}
	nPlayer.Password = hashedPassword

	id, err := h.playerRepository.Insert(ctx, &nPlayer)
	if err != nil {
		log.Printf("register insert player: %v", err)
		return nil, huma.Error500InternalServerError("create", err)
	}

	log.Printf("player %v created", id)
	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) SetPassword(
	ctx context.Context,
	i *RequestSetPassword,
) (*domain.ResponseID[string], error) {
	playerID, ok := ctxutil.GetPlayerID(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("player id is invalid")
	}

	nPlayer := player.PlayerUpdate{
		ID:       playerID,
		Password: &i.Body.Password,
	}
	if err := nPlayer.Valid(); err != nil {
		log.Printf("set pass not valid: %v", err)
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	hashedPassword, err := password.Hash(i.Body.Password)
	if err != nil {
		log.Printf("set pass hash: %v", err)
		return nil, huma.Error500InternalServerError("could not update player")
	}
	nPlayer.Password = &hashedPassword

	id, err := h.playerRepository.Update(ctx, &nPlayer)
	if err != nil {
		log.Printf("set pass player update: %v", err)
		return nil, huma.Error500InternalServerError("create", err)
	}

	log.Printf("player %v updated (pass)", id)
	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) issueToken(playerID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   playerID,
		ExpiresAt: jwt.NewNumericDate(now.Add(defaultExpiration)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}
