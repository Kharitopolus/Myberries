package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                  string
	DBUrl                    string
	RedisUrl                 string
	JWTSecret                string
	RefreshTokenExpTimeHours time.Duration
	AccessTokenExpTimeHours  time.Duration
}

func ConfigGet(filename string) Config {
	godotenv.Load(filename)
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatal("APP_PORT must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	return Config{
		AppPort: appPort,
		DBUrl:   dbURL,
	}
}
