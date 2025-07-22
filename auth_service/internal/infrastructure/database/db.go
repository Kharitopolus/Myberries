package database

import (
	"database/sql"
	"log"

	"github.com/Kharitopolus/Myberries/auth_service/internal/infrastructure/database/sqlc"
)

type DB struct {
	sqlc *sqlc.Queries
}

func New(url string) DB {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	return DB{sqlc: sqlc.New(conn)}
}
