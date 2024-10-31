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