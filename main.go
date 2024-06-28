package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"log"
	"net/http"
	"os"
)

import _ "github.com/lib/pq"

type apiConfig struct {
	DB        *database.Queries
	JWTSecret string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("$JWT_SECRET must be set")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		dbQueries,
		jwtSecret,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", apiCfg.handlerSignUp)
	mux.HandleFunc("POST /v1/users/login", apiCfg.handlerSignIn)
	mux.HandleFunc("POST /v1/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("GET /v1/users/me", apiCfg.middlewareAuth(apiCfg.handlerMe))
	mux.HandleFunc("POST /v1/songs", apiCfg.middlewareAuth(apiCfg.handlerSongsCreate))
	mux.HandleFunc("GET /v1/songs", apiCfg.middlewareAuth(apiCfg.handlerSongsGet))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting server on port %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
