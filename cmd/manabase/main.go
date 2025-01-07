package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/jsaterdalen/manabase"
	"github.com/jsaterdalen/manabase/cmd/web/views"
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

	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("cmd/web/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", indexHandler(queries))
	router.Get("/newgame", newGamePageHandler())
	router.Post("/game", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("error parsing form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)

		formData := r.Form
		for key, values := range formData {
			fmt.Print(key)
			fmt.Print("\n")
			for _, value := range values {
				fmt.Print(value)
				fmt.Print("\n")
			}
			fmt.Print("\n")
		}
	})

	srv := &http.Server{
		Handler: router,
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
