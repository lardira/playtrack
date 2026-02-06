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
	FindOne(ctx context.Context, id string) (*Player, error)
	Insert(context.Context, *Player) (string, error)
	Update(ctx context.Context, player *PlayerUpdate) (string, error)
}

type PlayedGameRepository interface {
	FindAll(ctx context.Context, playerID string) ([]PlayedGame, error)
	FindOne(ctx context.Context, id int) (*PlayedGame, error)
	FindLastNotReroll(ctx context.Context, playerID string) (*PlayedGame, error)
	Insert(ctx context.Context, game *PlayedGame) (int, error)
	Update(ctx context.Context, game *PlayedGame) (int, error)
}

type GameRepository interface {
	FindOne(ctx context.Context, id int) (*game.Game, error)
}

type Service struct {
	playerRepository     PlayerRepository
	playedGameRepository PlayedGameRepository
	gameService          *game.Service
}

func NewService(
	playerRepository PlayerRepository,
	playedGameRepository PlayedGameRepository,
	gameService *game.Service,
) *Service {
	return &Service{
		playerRepository:     playerRepository,
		playedGameRepository: playedGameRepository,
		gameService:          gameService,
	}
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
		return 0, err
	}

	id, err := s.playedGameRepository.Insert(ctx, nGamePlayed)
	if err != nil {
		log.Printf("played game insert: %v", err)
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdatePlayedGame(ctx context.Context, playerID string, playedGameID int, upd *PlayedGameUpdate) (int, error) {
	playedGame, err := s.GetOnePlayedGame(ctx, playedGameID)
	if err != nil {
		log.Printf("played find one: %v", err)
		return 0, err
	}

	if err := s.setUpdate(playedGame, upd); err != nil {
		return 0, err
	}
	if err := s.applyStatus(ctx, playerID, playedGame); err != nil {
		return 0, err
	}

	id, err := s.playedGameRepository.Update(ctx, playedGame)
	if err != nil {
		log.Printf("played game insert: %v", err)
		return 0, err
	}
	return id, nil
}

func (s *Service) setUpdate(pg *PlayedGame, upd *PlayedGameUpdate) error {
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

func (s *Service) applyStatus(ctx context.Context, playerID string, pg *PlayedGame) error {
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
