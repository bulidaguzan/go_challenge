package models

import "time"

// Transaction represents a financial transaction
type Transaction struct {
    ID       int64     `json:"id"`
    UserID   int64     `json:"user_id"`
    Amount   float64   `json:"amount"`
    DateTime time.Time `json:"datetime"`
}

// BalanceResponse represents the balance information response
type BalanceResponse struct {
    Balance      float64 `json:"balance" example:"25.21"`
    TotalDebits  int     `json:"total_debits" example:"10"`
    TotalCredits int     `json:"total_credits" example:"15"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error string `json:"error" example:"Invalid user ID"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
    Message string `json:"message" example:"Migration completed successfully"`
}

type MigrationStats struct {
    TotalRecords    int       `json:"total_records" example:"100"`
    SuccessfulRows  int       `json:"successful_rows" example:"95"`
    FailedRows      int       `json:"failed_rows" example:"5"`
    TotalAmount     float64   `json:"total_amount" example:"1500.50"`
    UniqueUsers     int       `json:"unique_users" example:"25"`
    DateRange       struct {
        Earliest    string    `json:"earliest" example:"2024-01-01T00:00:00Z"`
        Latest      string    `json:"latest" example:"2024-12-31T23:59:59Z"`
    } `json:"date_range"`
    TransactionTypes struct {
        Credits     int       `json:"credits" example:"60"`
        Debits      int       `json:"debits" example:"40"`
    } `json:"transaction_types"`
    Errors          []string  `json:"errors,omitempty" example:"Line 5: Invalid amount format"`
}