package player

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lardira/playtrack/internal/domain/game"
)

type PlayerRepository interface {
	FindAll(context.Context) ([]Player, error)
	FindOneByUsername(ctx context.Context, username string) (*Player, error)
	FindOne(ctx context.Context, id string) (*Player, error)
	Insert(context.Context, *Player) (string, error)
	Update(ctx context.Context, player *Player) (string, error)
}

type PlayedGameRepository interface {
	FindAll(ctx context.Context, playerID string) ([]PlayedGame, error)
	FindOne(ctx context.Context, id int) (*PlayedGame, error)
	FindLastNotReroll(ctx context.Context, playerID string) (*PlayedGame, error)
	Insert(ctx context.Context, game *PlayedGame) (int, error)
	Update(ctx context.Context, game *PlayedGame) (int, error)
}

type GameService interface {
	GetOne(ctx context.Context, gameID int) (*game.Game, error)
}

type Service struct {
	playerRepository     PlayerRepository
	playedGameRepository PlayedGameRepository
	gameService          GameService
}

func NewService(
	playerRepository PlayerRepository,
	playedGameRepository PlayedGameRepository,
	gameService GameService,
) *Service {
	return &Service{
		playerRepository:     playerRepository,
		playedGameRepository: playedGameRepository,
		gameService:          gameService,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Player, error) {
	players, err := s.playerRepository.FindAll(ctx)
	if err != nil {
		log.Printf("player find all: %v", err)
		return nil, err
	}

	return players, nil
}

func (s *Service) GetOne(ctx context.Context, playerID string) (*Player, error) {
	player, err := s.playerRepository.FindOne(ctx, playerID)
	if err != nil {
		log.Printf("player find one: %v", err)
		return nil, err
	}

	return player, nil
}

func (s *Service) GetOneByUsername(ctx context.Context, username string) (*Player, error) {
	player, err := s.playerRepository.FindOneByUsername(ctx, username)
	if err != nil {
		log.Printf("player find one: %v", err)
		return nil, err
	}

	return player, nil
}

func (s *Service) Create(ctx context.Context, params PlayerParams) (string, error) {
	nPlayer, err := NewPlayer(params)
	if err != nil {
		log.Printf("player create: %v", err)
		return "", err
	}

	id, err := s.playerRepository.Insert(ctx, nPlayer)
	if err != nil {
		log.Printf("insert player: %v", err)
		return "nil", err
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, playerID string, upd PlayerUpdate) (string, error) {
	player, err := s.GetOne(ctx, playerID)
	if err != nil {
		log.Printf("find one: %v", err)
		return "", err
	}

	if err := s.setUpdate(player, &upd); err != nil {
		log.Printf("set update params: %v", err)
		return "", err
	}

	id, err := s.playerRepository.Update(ctx, player)
	if err != nil {
		log.Printf("player update: %v", err)
		return "", err
	}
	return id, nil
}

func (s *Service) GetAllPlayedGames(ctx context.Context, playerID string) ([]PlayedGame, error) {
	games, err := s.playedGameRepository.FindAll(ctx, playerID)
	if err != nil {
		log.Printf("played games find all: %v", err)
		return nil, err
	}

	return games, nil
}

func (s *Service) GetOnePlayedGame(ctx context.Context, playedGameID int) (*PlayedGame, error) {
	game, err := s.playedGameRepository.FindOne(ctx, playedGameID)
	if err != nil {
		log.Printf("played games find one: %v", err)
		return nil, err
	}

	return game, nil
}

func (s *Service) CreatePlayedGame(ctx context.Context, playerID string, gameID int) (int, error) {
	game, err := s.gameService.GetOne(ctx, gameID)
	if err != nil {
		log.Printf("create played game: %v", err)
		return 0, err
	}

	// TODO: query db for this
	allPlayed, err := s.playedGameRepository.FindAll(ctx, playerID)
	if err != nil {
		return 0, err
	}
	for _, p := range allPlayed {
		if !p.StatusTerminated() {
			return 0, fmt.Errorf("player has game in nonterminated status: %v", p.ID)
		}
	}

	nGamePlayed, err := NewPlayedGame(PlayedGameParams{
		PlayerID:  playerID,
		GameID:    gameID,
		Points:    game.Points,
		StartedAt: time.Now(),
	})
	if err != nil {
		log.Printf("create played game: %v", err)
		return 0, err
	}

	id, err := s.playedGameRepository.Insert(ctx, nGamePlayed)
	if err != nil {
		log.Printf("played game insert: %v", err)
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdatePlayedGame(ctx context.Context, playerID string, playedGameID int, upd PlayedGameUpdate) (int, error) {
	playedGame, err := s.GetOnePlayedGame(ctx, playedGameID)
	if err != nil {
		log.Printf("update played game: %v", err)
		return 0, err
	}

	if err := s.setPlayedUpdate(playedGame, &upd); err != nil {
		log.Printf("update played game: %v", err)
		return 0, err
	}
	if err := s.applyPlayedStatus(ctx, playerID, playedGame); err != nil {
		log.Printf("update played game: %v", err)
		return 0, err
	}

	id, err := s.playedGameRepository.Update(ctx, playedGame)
	if err != nil {
		log.Printf("update played game: %v", err)
		return 0, err
	}
	return id, nil
}

func (s *Service) setUpdate(pg *Player, upd *PlayerUpdate) error {
	if upd.Username != nil {
		if err := pg.SetUsername(*upd.Username); err != nil {
			return err
		}
	}
	if upd.Img != nil {
		if err := pg.SetImage(*upd.Img); err != nil {
			return err
		}
	}
	if upd.Email != nil {
		if err := pg.SetEmail(*upd.Email); err != nil {
			return err
		}
	}
	if upd.Password != nil {
		if err := pg.SetPassword(*upd.Password); err != nil {
			return err
		}
	}
	if upd.Description != nil {
		if err := pg.SetDescription(*upd.Description); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) setPlayedUpdate(pg *PlayedGame, upd *PlayedGameUpdate) error {
	if upd.Points != nil {
		pg.Points = *upd.Points
	}
	if upd.Comment != nil {
		if err := pg.SetComment(*upd.Comment); err != nil {
			return err
		}
	}
	if upd.Rating != nil {
		if err := pg.SetRating(*upd.Rating); err != nil {
			return err
		}
	}
	if upd.StartedAt != nil {
		if err := pg.SetDates(*upd.StartedAt, upd.CompletedAt); err != nil {
			return err
		}
	} else if upd.CompletedAt != nil {
		return errors.New("completed date must be set with started date")
	}
	if upd.PlayTime != nil {
		pg.PlayTime = upd.PlayTime
	}
	if upd.Status != nil {
		if err := pg.SetStatus(*upd.Status); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) applyPlayedStatus(ctx context.Context, playerID string, pg *PlayedGame) error {
	switch pg.Status {
	case PlayedGameStatusDropped:
		dropPoints := -1

		prevGame, err := s.playedGameRepository.FindLastNotReroll(ctx, playerID)
		if err != nil && !errors.Is(err, ErrPlayedGameNotFound) {
			log.Printf("last played game find: %v", err)
			return err
		}
		// consecutive drops are stacked
		if err == nil && prevGame.Status == PlayedGameStatusDropped {
			dropPoints = prevGame.Points - 1
		}

		pg.Points = dropPoints
		if pg.CompletedAt == nil {
			if err := pg.Complete(); err != nil {
				return err
			}
		}

	case PlayedGameStatusRerolled:
		pg.Points = 0
		if pg.CompletedAt == nil {
			if err := pg.Complete(); err != nil {
				return err
			}
		}
	}
	return nil
}
