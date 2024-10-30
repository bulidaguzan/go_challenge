package routes

import (
	"database/sql"

	"github.com/bulidaguzan/go_challenge/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
    router := gin.Default()

    migrationHandler := handlers.NewMigrationHandler(db)
    balanceHandler := handlers.NewBalanceHandler(db)

    router.POST("/migrate", migrationHandler.MigrateCSV)
    router.GET("/users/:user_id/balance", balanceHandler.GetBalance)

    return router
}