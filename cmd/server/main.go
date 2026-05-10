package main

import (
	"log"
	"net/http"

	"expense_tracker_api/internal/handlers"
	"expense_tracker_api/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	repo, err := repository.NewExpenseRepository("expenses.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	expenseHandler := handlers.NewExpenseHandler(repo)

	r := chi.NewRouter()

	// Базовые middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Простой endpoint для проверки
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", expenseHandler.CreateExpense)
		r.Get("/", expenseHandler.GetExpenses)
		r.Get("/{id}", expenseHandler.GetExpense)
		r.Patch("/{id}", expenseHandler.UpdateExpense)
		r.Delete("/{id}", expenseHandler.DeleteExpense)
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
