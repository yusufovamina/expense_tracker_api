package main

import (
	"log"
	"net/http"
	"os"

	"expense_tracker_api/internal/handlers"
	"expense_tracker_api/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Игнорируем ошибку, если файла .env нет (будут использованы дефолтные значения)
	_ = godotenv.Load()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "expenses.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo, err := repository.NewExpenseRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	expenseHandler := handlers.NewExpenseHandler(repo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", expenseHandler.CreateExpense)
		r.Get("/", expenseHandler.GetExpenses)
		r.Get("/summary", expenseHandler.GetSummary)
		r.Get("/{id}", expenseHandler.GetExpense)
		r.Patch("/{id}", expenseHandler.UpdateExpense)
		r.Delete("/{id}", expenseHandler.DeleteExpense)
	})

	log.Println("Starting server on :" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
