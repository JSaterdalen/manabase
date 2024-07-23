package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jsaterdalen/manabase/cmd/web"
	"github.com/jsaterdalen/manabase/internal/database"

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

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	queries := database.New(conn)
	indexHandler := web.Handler{Queries: queries}

	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("cmd/web/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// router.Get("/", templ.Handler(web.HomePage()).ServeHTTP)
	router.Get("/", indexHandler.ServeHTTP)

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
