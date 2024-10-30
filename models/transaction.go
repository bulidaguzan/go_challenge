package models

import "time"

type Transaction struct {
    ID       int64
    UserID   int64
    Amount   float64
    DateTime time.Time
}

type BalanceResponse struct {
    Balance      float64 `json:"balance"`
    TotalDebits  int     `json:"total_debits"`
    TotalCredits int     `json:"total_credits"`
}