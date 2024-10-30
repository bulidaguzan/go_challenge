package db

import (
	"database/sql"
	"fmt"

	"github.com/bulidaguzan/go_challenge/config"
)

func InitDB(config config.Config) (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    if err := createSchema(db); err != nil {
        return nil, err
    }

    return db, nil
}

func createSchema(db *sql.DB) error {
    // Primero eliminar la tabla existente si existe
    _, err := db.Exec(`DROP TABLE IF EXISTS transactions;`)
    if err != nil {
        return fmt.Errorf("error dropping table: %w", err)
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
        return fmt.Errorf("error creating table: %w", err)
    }

    return nil
}
