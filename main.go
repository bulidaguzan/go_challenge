package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Config almacena la configuración de la aplicación
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    Port       string
}

// Transaction representa una transacción en la base de datos
type Transaction struct {
    ID       int64
    UserID   int64
    Amount   float64
    DateTime time.Time
}

// BalanceResponse representa la respuesta del endpoint de balance
type BalanceResponse struct {
    Balance      float64 `json:"balance"`
    TotalDebits  int     `json:"total_debits"`
    TotalCredits int     `json:"total_credits"`
}

// MigrationHandler maneja las operaciones de migración
type MigrationHandler struct {
    db *sql.DB
}

// NewMigrationHandler crea una nueva instancia de MigrationHandler
func NewMigrationHandler(db *sql.DB) *MigrationHandler {
    return &MigrationHandler{db: db}
}

// MigrateCSV maneja la solicitud de migración de CSV
func (h *MigrationHandler) MigrateCSV(c *gin.Context) {
    // Obtener el archivo del request
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }

    // Abrir el archivo
    openedFile, err := file.Open()
    if err != nil {
        c.JSON(500, gin.H{"error": "Error opening file"})
        return
    }
    defer openedFile.Close()

    // Crear el reader CSV
    reader := csv.NewReader(openedFile)

    // Leer la cabecera
    _, err = reader.Read()
    if err != nil {
        c.JSON(400, gin.H{"error": "Error reading CSV header"})
        return
    }

    // Preparar la consulta SQL que primero elimina y luego inserta
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

    // Leer y procesar cada línea
    for {
        record, err := reader.Read()
        if err != nil {
            break // Fin del archivo
        }

        // Convertir valores
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

        // Primero eliminar si existe y luego insertar
        _, err = stmt.Exec(id, userID, amount, datetime)
        if err != nil {
            log.Printf("Error processing record: %v", err)
            continue
        }
    }

    c.JSON(200, gin.H{"message": "Migration completed successfully"})
}


// BalanceHandler maneja las operaciones de balance
type BalanceHandler struct {
    db *sql.DB
}

// NewBalanceHandler crea una nueva instancia de BalanceHandler
func NewBalanceHandler(db *sql.DB) *BalanceHandler {
    return &BalanceHandler{db: db}
}

// GetBalance maneja la solicitud de obtención de balance
func (h *BalanceHandler) GetBalance(c *gin.Context) {
    // Obtener user_id del path
    userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid user ID"})
        return
    }

    // Obtener parámetros de fecha
    fromStr := c.Query("from")
    toStr := c.Query("to")

    // Construir la query base
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

    // Añadir filtros de fecha si están presentes
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

    // Ejecutar la query
    var response BalanceResponse
    err = h.db.QueryRow(query, args...).Scan(&response.Balance, &response.TotalDebits, &response.TotalCredits)
    if err != nil {
        c.JSON(500, gin.H{"error": "Database query error"})
        return
    }

    c.JSON(200, response)
}

// getConfig obtiene la configuración desde variables de entorno
func getConfig() Config {
    return Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "fintech"),
        Port:       getEnv("PORT", "8080"),
    }
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

// initDB inicializa la conexión a la base de datos
func initDB(config Config) (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    // Primero eliminar la tabla existente si existe
    _, err = db.Exec(`DROP TABLE IF EXISTS transactions;`)
    if err != nil {
        return nil, fmt.Errorf("error dropping table: %w", err)
    }

    // Crear la tabla exactamente como el CSV, sin constraints
    _, err = db.Exec(`
        CREATE TABLE transactions (
            id BIGINT,
            user_id BIGINT,
            amount DECIMAL(10,2),
            datetime TIMESTAMP WITH TIME ZONE
        );
        
        CREATE INDEX idx_transactions_user_id ON transactions(user_id);
        CREATE INDEX idx_transactions_datetime ON transactions(datetime);
    `)
    if err != nil {
        return nil, fmt.Errorf("error creating table: %w", err)
    }

    return db, nil
}

// setupRouter configura las rutas de la API
func setupRouter(db *sql.DB) *gin.Engine {
    router := gin.Default()

    migrationHandler := NewMigrationHandler(db)
    balanceHandler := NewBalanceHandler(db)

    router.POST("/migrate", migrationHandler.MigrateCSV)
    router.GET("/users/:user_id/balance", balanceHandler.GetBalance)

    return router
}

func main() {
    config := getConfig()

    db, err := initDB(config)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    router := setupRouter(db)

    log.Printf("Server starting on port %s", config.Port)
    if err := router.Run(":" + config.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}