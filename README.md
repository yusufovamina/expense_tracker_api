# Personal Expense Tracker API

A lightweight REST API for tracking personal expenses, built with Go and SQLite.

## How to Run

1. Make sure you have Go 1.21+ installed.
2. Clone the repository and navigate inside the folder.
3. Download dependencies:
   ```bash
   go mod tidy
   ```
4. Start the server (runs on `http://localhost:8080/`):
   ```bash
   go run ./cmd/server/main.go
   ```

A local `expenses.db` SQLite file will be automatically created in the root directory upon startup.

## Technology Choices

- **Go 1.23 & Chi Router (`go-chi/chi`)**: Chi is a lightweight, standard-library-compatible router. It provides the routing capabilities needed for REST APIs without the heavy footprint of larger frameworks.
- **SQLite3 (`mattn/go-sqlite3`)**: Perfect for a locally executed assignment. It demonstrates raw SQL usage, connection management, and persistence without requiring complex external database setups.

