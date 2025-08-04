package main

import (
	"github.com/Kharitopolus/Myberries/catalog_service/config"
	"github.com/Kharitopolus/Myberries/catalog_service/handlers"
)

func main() {
	cfg := config.ConfigGet(".env")

	db := handlers.InitDB(cfg.DBUrl)

	r := handlers.SetupRouter(db)

	r.Run()
}
