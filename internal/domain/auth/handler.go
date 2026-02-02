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
	found, err := h.playerRepository.FindOneByUsername(ctx, i.Body.Username)
	if err != nil {
		log.Printf("login player find one: %v", err)
		return nil, huma.Error401Unauthorized("username or password is incorrect")
	}
	if !password.CompareHash(i.Body.Password, found.Password) {
		log.Printf("login compare hash: %v", err)
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
	ctxPlr, ok := ctxutil.GetPlayer(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("player id is invalid")
	}

	found, err := h.playerRepository.FindOneByUsername(ctx, i.Body.Username)
	if err != nil {
		log.Printf("find one by username %v: %v", i.Body.Username, err)
		return nil, huma.Error401Unauthorized("player not found")
	}
	if !ctxPlr.IsAdmin && (found.ID != ctxPlr.ID) {
		log.Printf("player %v access to %v", ctxPlr, found.ID)
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	nPlayer := player.PlayerUpdate{
		ID:       found.ID,
		Username: &found.Username,
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
