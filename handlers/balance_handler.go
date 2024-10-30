// handlers/balance_handler.go
package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"fintech-backend/models"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
    db *sql.DB
}

func NewBalanceHandler(db *sql.DB) *BalanceHandler {
    return &BalanceHandler{db: db}
}

func (h *BalanceHandler) GetBalance(c *gin.Context) {
    userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid user ID"})
        return
    }

    fromStr := c.Query("from")
    toStr := c.Query("to")

    query := `
        SELECT 
            COALESCE(SUM(amount), 0) as balance,
            COUNT(CASE WHEN amount < 0 THEN 1 END) as debits,
            COUNT(CASE WHEN amount > 0 THEN 1 END) as credits
        FROM transactions
        WHERE user_id = $1
    `
    args := []interface{}{userID}
    argCount := 2

    if fromStr != "" && toStr != "" {
        from, err := time.Parse(time.RFC3339, fromStr)
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid from date format"})
            return
        }

        to, err := time.Parse(time.RFC3339, toStr)
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid to date format"})
            return
        }

        query += fmt.Sprintf(" AND datetime BETWEEN $%d AND $%d", argCount, argCount+1)
        args = append(args, from, to)
    }

    var response models.BalanceResponse
    err = h.db.QueryRow(query, args...).Scan(&response.Balance, &response.TotalDebits, &response.TotalCredits)
    if err != nil {
        c.JSON(500, gin.H{"error": "Database query error"})
        return
    }

    c.JSON(200, response)
}