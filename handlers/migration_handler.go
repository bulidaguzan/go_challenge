// handlers/migration_handler.go
package handlers

// MigrationHandler represents the migration handler
// @title Migration Handler
// @description Handler for migrating CSV data

import (
	"database/sql"
	"encoding/csv"
	"fintech-backend/models"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MigrationHandler struct {
    db *sql.DB
}


// MigrateCSV godoc
// @Summary Migrate CSV data
// @Description Upload and process a CSV file containing transaction records
// @Tags migration
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file to upload"
// @Success 200 {object} models.MigrationStats
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /migrate [post]
func NewMigrationHandler(db *sql.DB) *MigrationHandler {
    return &MigrationHandler{db: db}
}

func (h *MigrationHandler) MigrateCSV(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, models.ErrorResponse{Error: "No file uploaded"})
        return
    }

    openedFile, err := file.Open()
    if err != nil {
        c.JSON(500, models.ErrorResponse{Error: "Error opening file"})
        return
    }
    defer openedFile.Close()

    reader := csv.NewReader(openedFile)

    // Leer y validar headers
    headers, err := reader.Read()
    
    if err != nil {
        c.JSON(400, models.ErrorResponse{Error: "Error reading CSV header"})
        return
    }
    fmt.Print(headers)

    // Preparar estadísticas
    stats := models.MigrationStats{}
    uniqueUsers := make(map[int64]bool)
    var earliest, latest time.Time
    stats.DateRange.Earliest = "N/A"
    stats.DateRange.Latest = "N/A"

    stmt, err := h.db.Prepare(`
        WITH del AS (
            DELETE FROM transactions 
            WHERE id = $1
        )
        INSERT INTO transactions (id, user_id, amount, datetime)
        VALUES ($1, $2, $3, $4)
    `)
    if err != nil {
        c.JSON(500, models.ErrorResponse{Error: "Database preparation error"})
        return
    }
    defer stmt.Close()

    lineNumber := 1
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: Reading error", lineNumber))
            stats.FailedRows++
            continue
        }

        stats.TotalRecords++

        // Validar y procesar cada registro
        id, err := strconv.ParseInt(record[0], 10, 64)
        if err != nil {
            stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: Invalid ID format", lineNumber))
            stats.FailedRows++
            continue
        }

        userID, err := strconv.ParseInt(record[1], 10, 64)
        if err != nil {
            stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: Invalid user ID format", lineNumber))
            stats.FailedRows++
            continue
        }
        uniqueUsers[userID] = true

        amount, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: Invalid amount format", lineNumber))
            stats.FailedRows++
            continue
        }
        stats.TotalAmount += amount

        if amount > 0 {
            stats.TransactionTypes.Credits++
        } else {
            stats.TransactionTypes.Debits++
        }

        datetime, err := time.Parse(time.RFC3339, record[3])
        if err != nil {
            stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: Invalid datetime format", lineNumber))
            stats.FailedRows++
            continue
        }

        // Actualizar rango de fechas
        if stats.DateRange.Earliest == "N/A" || datetime.Before(earliest) {
            earliest = datetime
            stats.DateRange.Earliest = datetime.Format(time.RFC3339)
        }
        if stats.DateRange.Latest == "N/A" || datetime.After(latest) {
            latest = datetime
            stats.DateRange.Latest = datetime.Format(time.RFC3339)
        }

        _, err = stmt.Exec(id, userID, amount, datetime)
        if err != nil {
            stats.Errors = append(stats.Errors, fmt.Sprintf("Line %d: Database error", lineNumber))
            stats.FailedRows++
            continue
        }

        stats.SuccessfulRows++
        lineNumber++
    }

    stats.UniqueUsers = len(uniqueUsers)

    c.JSON(200, stats)
}
