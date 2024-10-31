package routes

import (
	"database/sql"
	"fintech-backend/handlers"

	"github.com/gin-gonic/gin"
)

// @title           Fintech Backend API
// @version         1.0
func SetupRouter(db *sql.DB) *gin.Engine {
    router := gin.Default()

    migrationHandler := handlers.NewMigrationHandler(db)
    balanceHandler := handlers.NewBalanceHandler(db)

    // @Summary      Upload CSV file
    // @Router       /migrate [post]
    router.POST("/migrate", migrationHandler.MigrateCSV)

    // @Summary      Get user balance
    // @Router       /users/{user_id}/balance [get]
    router.GET("/users/:user_id/balance", balanceHandler.GetBalance)

    return router
}