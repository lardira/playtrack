package player

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/domain"
	"github.com/lardira/playtrack/internal/domain/game"
	"github.com/lardira/playtrack/internal/pkg/ctxutil"
)

type PlayerRepository interface {
	FindAll(context.Context) ([]Player, error)
	FindOne(ctx context.Context, id string) (*Player, error)
	Insert(context.Context, *Player) (string, error)
	Update(ctx context.Context, player *PlayerUpdate) (string, error)
}

type PlayedGameRepository interface {
	FindAll(ctx context.Context, playerID string) ([]PlayedGame, error)
	FindOne(ctx context.Context, playerID string, id int) (*PlayedGame, error)
	FindLastNotReroll(ctx context.Context, playerID string) (*PlayedGame, error)
	Insert(ctx context.Context, player *PlayedGame) (int, error)
	Update(ctx context.Context, game *PlayedGameUpdate) (int, error)
}

type GameRepository interface {
	FindOne(ctx context.Context, id int) (*game.Game, error)
}

type Handler struct {
	playerRepository     PlayerRepository
	playedGameRepository PlayedGameRepository
	gameRepository       GameRepository
}

func NewHandler(
	playerRepository PlayerRepository,
	gameRepository GameRepository,
	playedGameRepository PlayedGameRepository,
) *Handler {
	return &Handler{
		playerRepository:     playerRepository,
		playedGameRepository: playedGameRepository,
		gameRepository:       gameRepository,
	}
}

func (h *Handler) Register(api huma.API) {
	grp := huma.NewGroup(api, "/players")
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"players"}
	})

	huma.Register(grp, huma.Operation{
		OperationID: "players-get-all",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "get all players",
		Description: "get all players",
	}, h.GetAll)

	huma.Register(grp, huma.Operation{
		OperationID: "players-get-one",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "get one player",
		Description: "get one player",
	}, h.GetOne)

	huma.Register(grp, huma.Operation{
		OperationID: "players-update-one",
		Method:      http.MethodPatch,
		Path:        "/{id}",
		Summary:     "update player",
		Description: "update a player",
	}, h.Update)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-get-all",
		Method:      http.MethodGet,
		Path:        "/{id}/played-games",
		Summary:     "get all played games",
		Description: "get all played games",
	}, h.GetAllPlayedGames)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-get-one",
		Method:      http.MethodGet,
		Path:        "/{id}/played-games/{gameID}",
		Summary:     "get one played game",
		Description: "get one played game",
	}, h.GetOnePlayedGame)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-create-one",
		Method:      http.MethodPost,
		Path:        "/{id}/played-games",
		Summary:     "create played game",
		Description: "create a new played game",
	}, h.CreatePlayedGame)

	huma.Register(grp, huma.Operation{
		OperationID: "played-games-update-one",
		Method:      http.MethodPatch,
		Path:        "/{id}/played-games/{gameID}",
		Summary:     "update played game",
		Description: "update a played game",
	}, h.UpdatePlayedGame)
}

func (h *Handler) GetAll(ctx context.Context, i *struct{}) (*domain.ResponseItems[Player], error) {
	players, err := h.playerRepository.FindAll(ctx)
	if err != nil {
		log.Printf("player find all: %v", err)
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[Player]{}
	resp.Body.Items = players
	return &resp, nil
}

func (h *Handler) GetOne(ctx context.Context, i *struct {
	ID string `path:"id" format:"uuid"`
}) (*domain.ResponseItem[Player], error) {
	player, err := h.playerRepository.FindOne(ctx, i.ID)
	if err != nil {
		log.Printf("player find one: %v", err)
		return nil, huma.Error500InternalServerError("find", err)
	}

	resp := domain.ResponseItem[Player]{}
	resp.Body.Item = player
	return &resp, nil
}

func (h *Handler) Update(
	ctx context.Context,
	i *RequestUpdatePlayer,
) (*domain.ResponseID[string], error) {
	if ok := checkAuthorizedFor(ctx, i.PlayerID); !ok {
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	nPlayer := PlayerUpdate{
		ID:          i.PlayerID,
		Username:    i.Body.Username,
		Img:         i.Body.Img,
		Email:       i.Body.Email,
		Description: i.Body.Description,
	}
	if err := nPlayer.Valid(); err != nil {
		log.Printf("player valid: %v", err)
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	id, err := h.playerRepository.Update(ctx, &nPlayer)
	if err != nil {
		log.Printf("player update: %v", err)
		return nil, huma.Error500InternalServerError("update", err)
	}

	log.Printf("player %v updated", id)
	resp := domain.ResponseID[string]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) GetAllPlayedGames(ctx context.Context, i *struct {
	ID string `path:"id" format:"uuid"`
}) (*domain.ResponseItems[PlayedGame], error) {
	games, err := h.playedGameRepository.FindAll(ctx, i.ID)
	if err != nil {
		log.Printf("played games find all: %v", err)
		return nil, huma.Error500InternalServerError("find all", err)
	}

	resp := domain.ResponseItems[PlayedGame]{}
	resp.Body.Items = games
	return &resp, nil
}

func (h *Handler) GetOnePlayedGame(ctx context.Context, i *struct {
	PlayerID string `path:"id" format:"uuid"`
	GameID   int    `path:"gameID"`
}) (*domain.ResponseItem[PlayedGame], error) {
	game, err := h.playedGameRepository.FindOne(ctx, i.PlayerID, i.GameID)
	if err != nil {
		log.Printf("played games find one: %v", err)
		return nil, huma.Error500InternalServerError("find", err)
	}

	resp := domain.ResponseItem[PlayedGame]{}
	resp.Body.Item = game
	return &resp, nil
}

func (h *Handler) CreatePlayedGame(
	ctx context.Context,
	i *RequestCreatePlayedGame,
) (*domain.ResponseID[int], error) {
	if ok := checkAuthorizedFor(ctx, i.PlayerID); !ok {
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	game, err := h.gameRepository.FindOne(ctx, i.Body.GameID)
	if err != nil {
		log.Printf("game find one: %v", err)
		return nil, huma.Error400BadRequest("game find", err)
	}

	nPlayed := PlayedGame{
		PlayerID: i.PlayerID,
		GameID:   i.Body.GameID,
		Points:   game.Points,
	}

	if err := nPlayed.Valid(); err != nil {
		log.Printf("played game valid: %v", err)
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	if err := h.containsNonterminatedPlayed(ctx, i.PlayerID); err != nil {
		log.Printf("player %v contains nonterminated: %v", i.PlayerID, err)
		return nil, err
	}

	id, err := h.playedGameRepository.Insert(ctx, &nPlayed)
	if err != nil {
		log.Printf("played game insert: %v", err)
		return nil, huma.Error500InternalServerError("create", err)
	}

	log.Printf("played game %v created", id)
	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) UpdatePlayedGame(
	ctx context.Context,
	i *RequestUpdatePlayedGame,
) (*domain.ResponseID[int], error) {
	if ok := checkAuthorizedFor(ctx, i.PlayerID); !ok {
		return nil, huma.Error403Forbidden("player cannot access this entity")
	}

	nGame := PlayedGameUpdate{
		ID:          i.GameID,
		Points:      i.Body.Points,
		Comment:     i.Body.Comment,
		Rating:      i.Body.Rating,
		Status:      i.Body.Status,
		CompletedAt: i.Body.CompletedAt,
		StartedAt:   i.Body.StartedAt,
		PlayTime:    i.Body.PlayTime,
	}
	if err := nGame.Valid(); err != nil {
		log.Printf("played game update valid: %v", err)
		return nil, huma.Error400BadRequest("entity is not valid", err)
	}

	playedGame, err := h.playedGameRepository.FindOne(ctx, i.PlayerID, i.GameID)
	if err != nil {
		log.Printf("played find one: %v", err)
		return nil, huma.Error400BadRequest("entity is not found", err)
	}

	if nGame.Status != nil {
		newStatus := *nGame.Status
		if err := playedGame.StatusNextValid(newStatus); err != nil {
			log.Printf("played game %v next status check: %v", playedGame.ID, err)
			return nil, huma.Error400BadRequest("entity is not valid", err)
		}

		newPoints := 0
		now := time.Now()

		switch newStatus {
		case PlayedGameStatusDropped:
			newPoints = -1

			prevGame, err := h.playedGameRepository.FindLastNotReroll(ctx, i.PlayerID)
			if err != nil && !errors.Is(err, ErrPlayedGameNotFound) {
				log.Printf("last played game find: %v", err)
				return nil, huma.Error400BadRequest("game played find", err)
			}

			// consecutive drops are stacked
			if err == nil && prevGame.Status == PlayedGameStatusDropped {
				newPoints = prevGame.Points - 1
			}

			nGame.Points = &newPoints
			if nGame.CompletedAt == nil {
				nGame.CompletedAt = &now
			}

		case PlayedGameStatusRerolled:
			nGame.Points = &newPoints
			if nGame.CompletedAt == nil {
				nGame.CompletedAt = &now
			}
		}
	}

	id, err := h.playedGameRepository.Update(ctx, &nGame)
	if err != nil {
		log.Printf("played game update: %v", err)
		return nil, huma.Error500InternalServerError("update", err)
	}

	log.Printf("played game %v updated", id)
	resp := domain.ResponseID[int]{}
	resp.Body.ID = id
	return &resp, nil
}

func (h *Handler) containsNonterminatedPlayed(ctx context.Context, playerID string) error {
	allPlayed, err := h.playedGameRepository.FindAll(ctx, playerID)
	if err != nil {
		return huma.Error400BadRequest("find played games", err)
	}

	for _, p := range allPlayed {
		if !p.StatusTerminated() {
			return huma.Error400BadRequest(
				fmt.Sprintf("player has game in nonterminated status: %v", p.ID),
			)
		}
	}
	return nil
}

func checkAuthorizedFor(ctx context.Context, playerID string) bool {
	ctxPlayer, ok := ctxutil.GetPlayer(ctx)
	if !ok {
		return false
	}

	if ctxPlayer.IsAdmin {
		return true
	} else if ctxPlayer.ID != playerID {
		return false
	}
	return true
}
