// main.go
package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/yourusername/project/config"
	"github.com/yourusername/project/db"
	"github.com/yourusername/project/routes"
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