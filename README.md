# Expense Tracker API

A simple REST API for tracking personal expenses built with Go and Beego v2.

## Overview

This project provides a lightweight expense tracker API with user registration, login, and expense management. Data is persisted in CSV files under the `data/` directory.

## Key Features

- User registration and authentication
- Expense creation, listing, update, and deletion
- Expense summary reporting
- Built-in health check endpoint
- Swagger/OpenAPI documentation available via the built-in Swagger UI

## Tech Stack

- Go 1.22
- Beego v2.1.0
- CSV-based persistence for users and expenses

## Repository Structure

- `main.go` - application entry point
- `routers/router.go` - route and namespace registration
- `controllers/` - request handlers and business logic
- `models/` - CSV persistence and domain models
- `data/` - user and expense CSV files (created automatically)
- `swagger/` - Swagger UI assets and generated swagger.json

## Prerequisites

- Go 1.22 installed
- `bee` CLI installed

### Install `bee`

If `bee` is not available after opening VS Code or creating a new terminal, install it with:

```bash
go install github.com/beego/bee/v2@latest
```

Then make sure your Go bin directory is on `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

If you want this to persist, add the export line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.).

## Run Locally

```bash
git clone <your-repo-url>
cd expense-tracker-api
go mod download

cp conf/app.conf.example conf/app.conf

go run main.go
```

The server will start on `http://localhost:8080`.

> Note: `conf/app.conf` is ignored in Git so each developer can keep local settings separately. The sample `conf/app.conf.example` is tracked and should be copied to `conf/app.conf` before running.

## Run with Bee CLI (recommended for docs)

```bash
bee run
```

If you want to auto-generate Swagger docs while running:

```bash
bee run -gendoc=true
```

## Available Endpoints

### Health

- `GET /api/v1/health`

### Authentication

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### Expenses

- `POST /api/v1/expenses`
- `GET /api/v1/expenses`
- `GET /api/v1/expenses/summary`
- `GET /api/v1/expenses/:id`
- `PUT /api/v1/expenses/:id`
- `DELETE /api/v1/expenses/:id`

### Expense list filters and sorting

The `GET /api/v1/expenses` endpoint supports query parameters for filtering and sorting:

- `category` — filter by category
- `date_from` — start date (YYYY-MM-DD)
- `date_to` — end date (YYYY-MM-DD)
- `sort_by` — `amount` or `expense_date`
- `sort_order` — `asc` or `desc`
- `limit` — max results

## Example Requests

### Register user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Fahim",
    "email": "fahim@gmail.com",
    "password": "123456"
  }'
```

### Login user

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "fahim@gmail.com",
    "password": "123456"
  }'
```

### Create expense

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "title": "Lunch",
    "amount": 12.50,
    "category": "Food",
    "note": "Office lunch",
    "expense_date": "2026-06-02"
  }'
```

## Swagger Documentation

The API Swagger UI is available at:

```text
http://localhost:8080/swagger/
```

To regenerate documentation:

```bash
bee generate docs
cp docs/swagger.json swagger/
```

## Tests

Run all unit tests with:

```bash
go test ./...
```

Check test coverage quickly:

```bash
go test ./... -cover
```

Generate a full coverage report and view it in your browser:

```bash
go test ./... -coverprofile=coverage.out

go tool cover -html=coverage.out
```

If `go tool cover` is not installed, install it with:

```bash
go install golang.org/x/tools/cmd/cover@latest
```

Then run:

```bash
go tool cover -html=coverage.out
```

## Data Storage

The application uses CSV files in the `data/` directory:

- `data/users.csv`
- `data/expenses.csv`

These files are created automatically when the app runs.

## Notes

- The app uses `Beego` namespaces to group routes under `/api/v1`
- The code is designed for local development and learning purposes
- Passwords are stored in plain text in the CSV data store, so do not use this for production systems
