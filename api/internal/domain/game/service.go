package game

import (
	"context"
	"log"
)

type GameRepository interface {
	FindAll(ctx context.Context) ([]Game, error)
	FindOne(ctx context.Context, id int) (*Game, error)
	Insert(ctx context.Context, game *Game) (int, error)
}

type Service struct {
	gameRepository GameRepository
}

func NewService(gameRepository GameRepository) *Service {
	return &Service{
		gameRepository: gameRepository,
	}
}

func (s *Service) CreateGame(ctx context.Context, req *RequestCreateGame) (int, error) {
	nGame, err := NewGame(req.Body.Title, req.Body.HoursToBeat, req.Body.URL)
	if err != nil {
		log.Printf("game validation: %v", err)
		return 0, err
	}

	id, err := s.gameRepository.Insert(ctx, nGame)
	if err != nil {
		log.Printf("game insert: %v", err)
		return 0, err
	}
	return id, nil
}

func (s *Service) GetOne(ctx context.Context, ID int) (*Game, error) {
	game, err := s.gameRepository.FindOne(ctx, ID)
	if err != nil {
		log.Printf("game find one: %v", err)
		return nil, err
	}
	return game, nil
}

func (s *Service) GetAll(ctx context.Context) ([]Game, error) {
	games, err := s.gameRepository.FindAll(ctx)
	if err != nil {
		log.Printf("game find all: %v", err)
		return nil, err
	}
	return games, nil
}
