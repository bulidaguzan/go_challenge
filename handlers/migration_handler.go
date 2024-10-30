// handlers/migration_handler.go
package handlers

import (
	"database/sql"
	"encoding/csv"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MigrationHandler struct {
    db *sql.DB
}

func NewMigrationHandler(db *sql.DB) *MigrationHandler {
    return &MigrationHandler{db: db}
}

func (h *MigrationHandler) MigrateCSV(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }

    openedFile, err := file.Open()
    if err != nil {
        c.JSON(500, gin.H{"error": "Error opening file"})
        return
    }
    defer openedFile.Close()

    reader := csv.NewReader(openedFile)

    _, err = reader.Read()
    if err != nil {
        c.JSON(400, gin.H{"error": "Error reading CSV header"})
        return
    }

    stmt, err := h.db.Prepare(`
        WITH del AS (
            DELETE FROM transactions 
            WHERE id = $1
        )
        INSERT INTO transactions (id, user_id, amount, datetime)
        VALUES ($1, $2, $3, $4)
    `)
    if err != nil {
        c.JSON(500, gin.H{"error": "Database preparation error"})
        return
    }
    defer stmt.Close()

    for {
        record, err := reader.Read()
        if err != nil {
            break
        }

        id, err := strconv.ParseInt(record[0], 10, 64)
        if err != nil {
            continue
        }

        userID, err := strconv.ParseInt(record[1], 10, 64)
        if err != nil {
            continue
        }

        amount, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            continue
        }

        datetime, err := time.Parse(time.RFC3339, record[3])
        if err != nil {
            continue
        }

        _, err = stmt.Exec(id, userID, amount, datetime)
        if err != nil {
            log.Printf("Error processing record: %v", err)
            continue
        }
    }

    c.JSON(200, gin.H{"message": "Migration completed successfully"})
}
