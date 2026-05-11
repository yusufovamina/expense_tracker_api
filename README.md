# Personal Expense Tracker API

A lightweight REST API for a personal expense tracking application.

## How to Run

1. Clone the repository and ensure you have Go 1.21+ installed.
2. Download dependencies:
   ```bash
   go mod tidy
   ```
3. Start the server:
   ```bash
   go run ./cmd/server/main.go
   ```
*(A local SQLite database `expenses.db` will be automatically generated upon the first run).*

## Technology Choices

- **Go & Chi Router:** I chose `chi` because it is idiomatic, relies squarely on standard `net/http` handlers, and adds exactly what is needed (routing and basic middleware) without the heavy footprint of larger frameworks. 
- **SQLite:** Chosen for data persistence because it requires zero external setups or database containers. It acts as a perfect vehicle to demonstrate raw SQL querying, aggregations, and repository patterns effectively while remaining extremely easy to review locally

## What I Would Improve

If I had more time, I would proactively add:
1. **Testing:** Covering the core domain with table-driven unit tests and an integration test suite that spins up a mock DB.
2. **Pagination:** Implementing `LIMIT` and `OFFSET` on the `GET /expenses` endpoint. Returning the entire table works locally but will eventually become a performance bottleneck.
3. **Structured Logging:** Swapping the standard library `log` for `slog` to provide structured JSON logs, which are much easier to parse and monitor in real production environments.

## Assumptions Made

- **PATCH endpoint behavior:** The spec required an endpoint to update "amount, category, or note". I assumed this meant true partial updates. My implementation parses the JSON payload directly over the existing database record, leaving omitted fields completely untouched, while still enforcing the strict `amount > 0` validation.
- **Category data type:** I assumed categories are unstructured free-text strings rather than a strict enum. This gives the API consumer maximum flexibility to create tags on the fly.
- **Spent_on date mapping:** Since the sample JSON strictly used `"2024-03-10"`, I assumed storing it natively as a `YYYY-MM-DD` string natively maps effectively to SQLite's `TEXT` column, eliminating heavy `time.Time` unmarshaling overhead while still natively supporting SQL sorting and aggregations.
