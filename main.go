package main

import (
	"log"

	"fintech-backend/config"
	"fintech-backend/db"
	"fintech-backend/routes"

	_ "fintech-backend/docs"

	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Fintech Backend API
// @version         1.0
// @description     API for managing financial transactions and user balances
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @basePath  /
// @schemes   http https



func main() {
    config := config.GetConfig()

    db, err := db.InitDB(config)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    router := routes.SetupRouter(db)

    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Printf("Server starting on port %s", config.Port)
    if err := router.Run(":" + config.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}