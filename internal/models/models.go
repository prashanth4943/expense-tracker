package models

// Expense stores amount in minor units (paise) to avoid floating-point issues.
// e.g. ₹10.50 is stored as 1050.
type Expense struct {
	ID             string `json:"id"`
	AmountPaise    int64  `json:"amount_paise"` // stored value
	AmountDisplay  string `json:"amount"`       // formatted ₹ string for clients
	Category       string `json:"category"`
	Description    string `json:"description"`
	Date           string `json:"date"`       // YYYY-MM-DD
	CreatedAt      string `json:"created_at"` // RFC3339
	IdempotencyKey string `json:"-"`          // not exposed to client
}

type CreateExpenseRequest struct {
	Amount         string `json:"amount"` // e.g. "10.50"
	Category       string `json:"category"`
	Description    string `json:"description"`
	Date           string `json:"date"`            // YYYY-MM-DD
	IdempotencyKey string `json:"idempotency_key"` // client-generated UUID
}

type ListExpensesResponse struct {
	Expenses   []Expense `json:"expenses"`
	TotalPaise int64     `json:"total_paise"`
	Total      string    `json:"total"` // formatted
}
