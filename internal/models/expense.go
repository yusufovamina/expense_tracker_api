package models

import "time"

type Expense struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	Note      string    `json:"note,omitempty"`
	SpentOn   string    `json:"spent_on"`
	CreatedAt time.Time `json:"created_at"`
}

type Summary struct {
	TotalAmount float64            `json:"total_amount"`
	ByCategory  map[string]float64 `json:"by_category"`
}
