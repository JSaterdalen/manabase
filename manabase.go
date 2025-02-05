package manabase

import (
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
	Owner     Player
}

type Game struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DatePlayed time.Time
	GameNumber int32
	IsTotem    bool
	Players    []GamePlayer
	Winner     GamePlayer
}

type GamePlayer struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Deck      Deck
}

type Player struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (p GamePlayer) IsWinner(g Game) bool {
	return p.ID == g.Winner.ID
}

// type PlayerService interface {
// 	GetPlayerList() ([]GamePlayer, error)
// }

// type playerService struct {
// 	queries *database.Queries
// 	context context.Context
// }

// func NewPlayerService(context context.Context, queries *database.Queries) PlayerService {
// 	return &playerService{
// 		queries: queries,
// 		context: context,
// 	}
// }

// func (s *playerService) GetPlayerList() ([]GamePlayer, error) {
// 	rows, err := s.queries.GetPlayers(s.context)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return MakePlayers(rows), nil
// }

func MakeGames(gameRow []database.GetPlayerDeckGameRow) (games []Game) {
	gameMap := make(map[string]Game)
	for _, game := range gameRow {
		gameId := game.GameID.String()
		if _, ok := gameMap[gameId]; !ok {
			gameMap[gameId] = Game{
				ID:         game.GameID,
				DatePlayed: game.DatePlayed,
				Players:    []GamePlayer{},
				Winner:     GamePlayer{},
			}
		}

		player := GamePlayer{
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
		if games[i].DatePlayed.Equal(games[j].DatePlayed) {
			return games[i].ID.String() < games[j].ID.String()
		}
		return games[i].DatePlayed.After(games[j].DatePlayed)
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
