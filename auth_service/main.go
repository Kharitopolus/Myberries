package main

import (
	"log"
	"net/http"

	"github.com/Kharitopolus/Myberries/auth_service/internal/config"
	"github.com/Kharitopolus/Myberries/auth_service/internal/handlers"
	"github.com/Kharitopolus/Myberries/auth_service/internal/infrastructure/auth"
	"github.com/Kharitopolus/Myberries/auth_service/internal/infrastructure/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Get(".env")

	mux := http.NewServeMux()

	db := database.New(cfg.DBUrl)
	pm := auth.PasswordManager{}
	tm := auth.TokenManager{
		TokenSecret:         cfg.JWTSecret,
		AccessTokenExpTime:  cfg.AccessTokenExpTimeHours,
		RefreshTokenExpTime: cfg.RefreshTokenExpTimeHours,
	}

	uh := handlers.NewUsersHandlersImpl(db, pm, tm)

	uh.Router(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", cfg.AppPort)
	log.Fatal(srv.ListenAndServe())
}
