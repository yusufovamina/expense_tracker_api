package handlers

import (
	"encoding/json"
	"net/http"

	"expense_tracker_api/internal/models"
	"expense_tracker_api/internal/repository"
)

type ExpenseHandler struct {
	repo *repository.ExpenseRepository
}

func NewExpenseHandler(repo *repository.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{repo: repo}
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var e models.Expense
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if e.Amount <= 0 {
		http.Error(w, "Amount must be a positive number greater than 0", http.StatusBadRequest)
		return
	}
	if e.Category == "" || e.SpentOn == "" {
		http.Error(w, "Category and spent_on are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&e); err != nil {
		http.Error(w, "Failed to create expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch expenses", http.StatusInternalServerError)
		return
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
