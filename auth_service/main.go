package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Kharitopolus/Myberries/auth_service/internal/database"
	"github.com/Kharitopolus/Myberries/auth_service/internal/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db        *database.Queries
	rd        redis.RClient
	platform  string
	jwtSecret string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL must be set")
	}
	refreshTokenExpTimeHoursStr := os.Getenv("REFRESH_TOKEN_EXP_TIME_HOURS")
	if refreshTokenExpTimeHoursStr == "" {
		log.Fatal("REFRESH_TOKEN_EXP_TIME_HOURS must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	refreshToeknExpTimeHours, err := strconv.Atoi(refreshTokenExpTimeHoursStr)
	if err != nil {
		log.Fatal("invalid value in REFRESH_TOKEN_EXPT_TIME_HOURS")
	}

	rClient := redis.New(redisURL, refreshToeknExpTimeHours)

	apiCfg := apiConfig{
		db:        dbQueries,
		rd:        rClient,
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", apiCfg.handlerRegister)
	mux.HandleFunc("POST /auth/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /auth/me", apiCfg.handlerMe)
	mux.HandleFunc("POST /auth/refresh", apiCfg.handlerRefresh)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
