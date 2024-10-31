package db

import (
	"database/sql"
	"fmt"

	"fintech-backend/config"

	_ "github.com/lib/pq"
)

// Variable para facilitar el mockeo en tests
var sqlOpen = sql.Open

func InitDB(config config.Config) (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

    db, err := sqlOpen("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    if err := createSchema(db); err != nil {
        return nil, fmt.Errorf("error creating schema: %w", err)
    }

    return db, nil
}

func createSchema(db *sql.DB) error {
    // Drop table if exists
    _, err := db.Exec(`DROP TABLE IF EXISTS transactions;`)
    if err != nil {
        return fmt.Errorf("error dropping table: %w", err)
    }

    // Create table and indexes
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
        return fmt.Errorf("error creating table: %w", err)
    }

    return nil
}