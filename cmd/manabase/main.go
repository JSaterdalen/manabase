package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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
	// indexHandler := views.Handler{Queries: queries}

	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// router.Get("/", )
	router.Get("/", testHandler())

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

func testHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var name string
		// Execute the query.
		row := db.QueryRow("SELECT myname FROM mytable")
		if err := row.Scan(&name); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Write it back to the client.
		fmt.Fprintf(w, "hi %s!\n", name)
	})
}
