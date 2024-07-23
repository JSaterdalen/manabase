package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jsaterdalen/manabase/cmd/web"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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

	// conn, err := sql.Open("postgres", dbUrl)
	// if err != nil {
	// 	log.Fatal("Can't connect to database", err)
	// }

	// queries := database.New(conn)

	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("cmd/web/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", templ.Handler(web.HomePage()).ServeHTTP)
	// router.Get("/games", func(w http.ResponseWriter, r *http.Request) {
	// 	players := []manabase.Player{
	// 		{UUID: "1", Name: "John"},
	// 		{UUID: "2", Name: "Jane"},
	// 	}
	// 	games := []manabase.Game{
	// 		{UUID: "1", Players: players, DatePlayed: "2021-01-01", GameNumber: 1},
	// 	}

	// 	component := web.GameList(games)
	// 	component.Render(r.Context(), w)
	// })

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server listening on port %s", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
