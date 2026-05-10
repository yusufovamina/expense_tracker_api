package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"expense_tracker_api/internal/models"
	"expense_tracker_api/internal/repository"

	"github.com/go-chi/chi/v5"
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

func (h *ExpenseHandler) GetExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	expense, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch expense", http.StatusInternalServerError)
		return
	}
	if expense == nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	expense, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch expense", http.StatusInternalServerError)
		return
	}
	if expense == nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(expense); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if expense.Amount <= 0 {
		http.Error(w, "Amount must be a positive number greater than 0", http.StatusBadRequest)
		return
	}
	if expense.Category == "" || expense.SpentOn == "" {
		http.Error(w, "Category and spent_on are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(id, expense); err != nil {
		http.Error(w, "Failed to update expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	found, err := h.repo.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
