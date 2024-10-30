package main

import (
	"log"

	"fintech-backend/config"
	"fintech-backend/db"
	"fintech-backend/routes"

	_ "github.com/lib/pq"
)

func main() {
    config := config.GetConfig()

    db, err := db.InitDB(config)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    router := routes.SetupRouter(db)

    log.Printf("Server starting on port %s", config.Port)
    if err := router.Run(":" + config.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}