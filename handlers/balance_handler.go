package handlers

import (
	"database/sql"
	"fintech-backend/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
    db *sql.DB
}

// GetBalance godoc
// @Summary Get user balance
// @Description Get the balance and transaction counts for a specific user
// @Tags balance
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param from query string false "Start date (RFC3339) 2024-07-05T00:00:00"
// @Param to query string false "End date (RFC3339) 2024-07-05T00:00:00"
// @Success 200 {object} models.BalanceResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{user_id}/balance [get]
func NewBalanceHandler(db *sql.DB) *BalanceHandler {
    return &BalanceHandler{db: db}
}

func (h *BalanceHandler) GetBalance(c *gin.Context) {
    userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
    if err != nil {
        c.JSON(400, models.ErrorResponse{Error: "Invalid user ID"})
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
            c.JSON(400, models.ErrorResponse{Error: "Invalid from date format"})
            return
        }

        to, err := time.Parse(time.RFC3339, toStr)
        if err != nil {
            c.JSON(400, models.ErrorResponse{Error: "Invalid to date format"})
            return
        }

        query += fmt.Sprintf(" AND datetime BETWEEN $%d AND $%d", argCount, argCount+1)
        args = append(args, from, to)
    }

    var response models.BalanceResponse
    err = h.db.QueryRow(query, args...).Scan(&response.Balance, &response.TotalDebits, &response.TotalCredits)
    if err != nil {
        c.JSON(500, models.ErrorResponse{Error: "Database query error"})
        return
    }

    c.JSON(200, response)
}