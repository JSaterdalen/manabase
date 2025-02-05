package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jsaterdalen/manabase"
	"github.com/jsaterdalen/manabase/cmd/web/views"
	"github.com/jsaterdalen/manabase/cmd/web/views/components"
	"github.com/jsaterdalen/manabase/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		fmt.Println("PORT is not found in the environment")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		fmt.Println("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	queries := database.New(conn)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("cmd/web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/newgame", newGamePageHandler())
	mux.HandleFunc("/newgame/players", newGamePlayers(queries))
	mux.HandleFunc("/newgame/decks", newGameDeck(queries))
	mux.HandleFunc("/game", gameHandler(queries))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			indexHandler(queries)(w, r)
			return
		}

		http.NotFound(w, r)
	})

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + portString,
	}

	log.Printf("Server listening on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := queries.GetPlayerDeckGame(r.Context())
		if err != nil {
			fmt.Println("Error getting games", err)
		}

		games := manabase.MakeGames(rows)

		views.HomePage(games).Render(r.Context(), w)
	}
}

func newGamePageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.NewGamePage().Render(r.Context(), w)
	}
}

func gameHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			log.Printf("error parsing form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the date from the form
		date := r.Form.Get("date")
		if date == "" {
			components.DateValidation().Render(r.Context(), w)
		}

		// Parse player-deck selections
		playerDecks := parsePlayerDecks(r.Form)

		// Create the game
		gameID := uuid.New()
		err = queries.InsertGame(r.Context(), database.InsertGameParams{
			ID:         gameID,
			DatePlayed: time.Now(),
			UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			CreatedAt:  time.Now(),
		})
		if err != nil {
			log.Printf("error creating game: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		winnerID := uuid.MustParse(r.Form.Get("winner"))

		// Create player_deck_game entries
		for playerID, deckID := range playerDecks {
			err = queries.InsertPlayerDeckGames(r.Context(), database.InsertPlayerDeckGamesParams{
				ID:        uuid.New(),
				PlayerID:  playerID,
				DeckID:    deckID,
				GameID:    gameID,
				CreatedAt: time.Now(),
				UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
				IsWon:     winnerID == playerID,
			})

			if err != nil {
				log.Printf("error creating player_deck_game: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func newGamePlayers(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := queries.GetPlayers(r.Context())
		if err != nil {
			fmt.Println("Error getting players", err)
		}

		players := manabase.MakePlayers(rows)

		components.PlayerFields(players).Render(r.Context(), w)
	}
}

func newGameDeck(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			log.Printf("error parsing form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get all player information and sort by name
		players := make([]manabase.Player, 0)
		for _, playerID := range r.Form["players"] {
			id := uuid.MustParse(playerID)
			player, err := queries.GetPlayer(r.Context(), id)
			if err != nil {
				fmt.Println("Error getting player", err)
				continue
			}
			players = append(players, manabase.Player{
				ID:   player.ID,
				Name: player.Name,
			})
		}

		// Sort players by name
		sort.Slice(players, func(i, j int) bool {
			return strings.ToLower(players[i].Name) < strings.ToLower(players[j].Name)
		})

		// Create sorted playerMap and playerIDs
		playerIDs := make([]uuid.UUID, 0)
		playerMap := make(map[uuid.UUID]string)
		for _, p := range players {
			playerIDs = append(playerIDs, p.ID)
			playerMap[p.ID] = p.Name
		}

		// Get all decks for each player, sorted by when that player last played them
		decksByPlayer := make(map[uuid.UUID][]manabase.Deck)
		for _, playerID := range playerIDs {
			rows, err := queries.GetDecksByLastPlayedByPlayer(r.Context(), database.GetDecksByLastPlayedByPlayerParams{
				Player:    playerID,  // All decks from selected players
				Playerids: playerIDs, // Sort by this player's play history
			})
			if err != nil {
				fmt.Println("Error getting decks for player", playerID, err)
				continue
			}

			decks := make([]manabase.Deck, 0)
			for _, row := range rows {
				deck := manabase.Deck{
					ID:   row.ID,
					Name: row.Name,
				}
				decks = append(decks, deck)
			}
			decksByPlayer[playerID] = decks
		}

		components.DeckFields(playerMap, decksByPlayer).Render(r.Context(), w)
	}
}

// parsePlayerDecks extracts player-deck pairs from form data
// Form fields should be in the format deck-[player-uuid]
func parsePlayerDecks(form map[string][]string) map[uuid.UUID]uuid.UUID {
	playerDecks := make(map[uuid.UUID]uuid.UUID)
	for key, values := range form {
		if len(values) == 0 {
			continue
		}

		// Check if this is a deck field
		if !strings.HasPrefix(key, "deck-") {
			continue
		}

		// Extract player UUID from the key
		playerIDStr := strings.TrimPrefix(key, "deck-")
		playerID, err := uuid.Parse(playerIDStr)
		if err != nil {
			continue
		}

		// Parse deck UUID from value
		deckID, err := uuid.Parse(values[0])
		if err != nil {
			continue
		}

		playerDecks[playerID] = deckID
	}
	return playerDecks
}
