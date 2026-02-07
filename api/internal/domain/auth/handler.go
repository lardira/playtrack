package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/domain/player"
	"github.com/lardira/playtrack/internal/middleware"
	"github.com/lardira/playtrack/internal/pkg/apiutil"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
)

const (
	defaultExpiration = 10 * 24 * time.Hour
)

type Handler struct {
	secret        string
	playerService *player.Service
}

func NewHandler(secret string, playerService *player.Service) *Handler {
	return &Handler{
		secret:        secret,
		playerService: playerService,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/auth")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"auth"}
	})

	huma.Register(grp, huma.Operation{
		OperationID: "register-player",
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "register player",
		Description: "register a new player (auth and player entity will be created)",
	}, h.RegisterPlayer)

	huma.Register(grp, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "login",
		Description: "login a player in order to get a jwt token",
	}, h.Login)

	huma.Register(grp, huma.Operation{
		OperationID: "set-password",
		Method:      http.MethodPatch,
		Path:        "/set-password",
		Summary:     "set new password",
		Description: "set new password for a player secured",
		Security:    apiutil.OperationSecurity,
		Middlewares: huma.Middlewares{middleware.Authorize(h.secret)},
	}, h.SetPassword)
}

func (h *Handler) Login(ctx context.Context, i *RequestLoginPlayer) (*ResponseLoginPlayer, error) {
	found, err := h.playerService.GetOneByUsername(ctx, i.Body.Username)
	if err != nil {
		return nil, huma.Error401Unauthorized("username or password is incorrect")
	}
	if !found.CheckPassword(i.Body.Password) {
		log.Printf("login compare hash error for user %s", i.Body.Username)
		return nil, huma.Error401Unauthorized("username or password is incorrect")
	}

	token, err := h.issueToken(found)
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
	i *player.RequestCreatePlayer,
) (*domain.ResponseID[string], error) {
	id, err := h.playerService.Create(ctx, player.PlayerParams{
		Username: i.Body.Username,
		Password: i.Body.Password,
		Img:      i.Body.Img,
		Email:    i.Body.Email,
	})
	if err != nil {
		return nil, huma.Error400BadRequest("register:", err)
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
	ctxPlr, ok := ctxutil.GetPlayer(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("player id is invalid")
	}

	p, err := h.playerService.GetOne(ctx, ctxPlr.ID)
	if err != nil {
		return nil, huma.Error404NotFound("player not found", err)
	}

	if !ctxPlr.IsAdmin && (p.Username != i.Body.Username) {
		log.Printf("player %v access to %v", ctxPlr, p.Username)
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	id, err := h.playerService.Update(ctx, p.ID, player.PlayerUpdate{
		Password: &i.Body.Password,
	})
	if err != nil {
		log.Printf("set pass player update: %v", err)
		return nil, huma.Error500InternalServerError("set password", err)
	}

	log.Printf("player %v updated (pass)", id)
	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) issueToken(p *player.Player) (string, error) {
	now := time.Now()
	audience := []string{apiutil.RolePlayer}

	if p.IsAdmin {
		audience = append(audience, apiutil.RoleAdmin)
	}

	claims := jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   p.ID,
		ExpiresAt: jwt.NewNumericDate(now.Add(defaultExpiration)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		Audience:  audience,
	}

	token := jwt.NewWithClaims(apiutil.DefaultSigningMethod, claims)
	return token.SignedString([]byte(h.secret))
}
