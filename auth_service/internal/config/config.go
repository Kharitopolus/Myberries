package config

import (
	"log"
	"os"
	"strconv"
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

func Get() Config {
	godotenv.Load()
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatal("APP_PORT must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL must be set")
	}
	accessTokenExpTimeHoursStr := os.Getenv("ACCESS_TOKEN_EXP_TIME_HOURS")
	if accessTokenExpTimeHoursStr == "" {
		log.Fatal("ACCESS_TOKEN_EXP_TIME_HOURS must be set")
	}
	refreshTokenExpTimeHoursStr := os.Getenv("REFRESH_TOKEN_EXP_TIME_HOURS")
	if refreshTokenExpTimeHoursStr == "" {
		log.Fatal("REFRESH_TOKEN_EXP_TIME_HOURS must be set")
	}
	acceessTokenExpTimeHours, err := strconv.Atoi(accessTokenExpTimeHoursStr)
	if err != nil {
		log.Fatal("invalid value in REFRESH_TOKEN_EXPT_TIME_HOURS")
	}
	refreshToeknExpTimeHours, err := strconv.Atoi(refreshTokenExpTimeHoursStr)
	if err != nil {
		log.Fatal("invalid value in REFRESH_TOKEN_EXPT_TIME_HOURS")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	return Config{
		AppPort:   appPort,
		DBUrl:     dbURL,
		RedisUrl:  redisURL,
		JWTSecret: jwtSecret,
		AccessTokenExpTimeHours: time.Duration(
			acceessTokenExpTimeHours,
		) * time.Hour,
		RefreshTokenExpTimeHours: time.Duration(
			refreshToeknExpTimeHours,
		) * time.Hour,
	}
}
