package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/prashanth4943/expense-tracker/internal/db"
	"github.com/prashanth4943/expense-tracker/internal/models"
	"github.com/prashanth4943/expense-tracker/internal/utils"
)

type Handler struct {
	store *db.Store
}

func NewHandler(store *db.Store) *Handler {
	return &Handler{store: store}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req models.CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// --- Validation ---
	req.Category = strings.TrimSpace(req.Category)
	req.Description = strings.TrimSpace(req.Description)
	req.Date = strings.TrimSpace(req.Date)
	req.IdempotencyKey = strings.TrimSpace(req.IdempotencyKey)

	if req.IdempotencyKey == "" {
		writeError(w, http.StatusBadRequest, "idempotency_key is required")
		return
	}
	if req.Category == "" {
		writeError(w, http.StatusBadRequest, "category is required")
		return
	}
	if req.Date == "" {
		writeError(w, http.StatusBadRequest, "date is required")
		return
	}

	amountPaise, err := utils.ParseAmount(req.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, "amount: "+err.Error())
		return
	}

	expense, err := h.store.CreateExpense(req, amountPaise)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create expense")
		return
	}

	writeJSON(w, http.StatusCreated, expense)
}

// ListExpenses handles GET /expenses
func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := db.ListFilter{
		Category: strings.TrimSpace(q.Get("category")),
		Sort:     q.Get("sort"),
	}

	expenses, err := h.store.ListExpenses(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list expenses")
		return
	}

	// Compute total of visible expenses.
	var total int64
	for _, e := range expenses {
		total += e.AmountPaise
	}

	if expenses == nil {
		expenses = []models.Expense{} // return [] not null
	}

	writeJSON(w, http.StatusOK, models.ListExpensesResponse{
		Expenses:   expenses,
		TotalPaise: total,
		Total:      utils.FormatAmount(total),
	})
}

// ListCategories handles GET /categories
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.store.Categories()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}
	if cats == nil {
		cats = []string{}
	}
	writeJSON(w, http.StatusOK, cats)
}
