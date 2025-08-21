package main

import (
	"log"
	"os"

	"ergracer-api/internal/api"
	"ergracer-api/internal/config"
	"ergracer-api/internal/database"
)

func main() {
	cfg := config.Load()
	log.Printf("the db URL: %s", cfg.DatabaseURL())
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	server := api.NewServer(db, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
