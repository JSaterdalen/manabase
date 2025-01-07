package manabase

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jsaterdalen/manabase/internal/database"
)

type Deck struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Commander string
	OwnerID   uuid.UUID
}

type Game struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DatePlayed time.Time
	GameNumber int32
	IsTotem    bool
	Players    []Player
	Winner     Player
}

type Player struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Deck      Deck
}

func (p Player) IsWinner(g Game) bool {
	return p.ID == g.Winner.ID
}

type PlayerService interface {
	GetPlayerList() ([]Player, error)
}

type playerService struct {
	queries *database.Queries
	context context.Context
}

func NewPlayerService(context context.Context, queries *database.Queries) PlayerService {
	return &playerService{
		queries: queries,
		context: context,
	}
}

func (s *playerService) GetPlayerList() ([]Player, error) {
	rows, err := s.queries.GetPlayers(s.context)
	if err != nil {
		return nil, err
	}

	return MakePlayers(rows), nil
}

func MakeGames(gameRow []database.GetPlayerDeckGameRow) (games []Game) {
	gameMap := make(map[string]Game)
	for _, game := range gameRow {
		gameId := game.GameID.String()
		if _, ok := gameMap[gameId]; !ok {
			gameMap[gameId] = Game{
				ID:         game.GameID,
				GameNumber: game.GameNumber,
				DatePlayed: game.DatePlayed,
				Players:    []Player{},
				Winner:     Player{},
			}
		}

		player := Player{
			ID:   game.PlayerID,
			Name: game.PlayerName,
			Deck: Deck{
				ID:   game.DeckID,
				Name: game.DeckName,
			},
		}
		g := gameMap[gameId]
		g.Players = append(g.Players, player)
		sort.Slice(g.Players, func(i, j int) bool {
			return g.Players[i].Name < g.Players[j].Name
		})

		if game.IsWon {
			g.Winner = player
		}
		gameMap[gameId] = g
	}

	for _, value := range gameMap {
		games = append(games, value)
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].GameNumber > games[j].GameNumber
	})
	return games
}

func MakePlayers(playerRows []database.Player) (players []Player) {
	for _, player := range playerRows {
		players = append(players, Player{
			ID:   player.ID,
			Name: player.Name,
		})
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})
	return players
}
