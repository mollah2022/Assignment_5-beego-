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
git clone https://github.com/mollah2022/Assignment_5-beego-.git
cd Assignment_5-beego-
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
  - No authentication required

### Authentication

- `POST /api/v1/auth/register`
  - Required JSON body:
    - `name` (string)
    - `email` (string)
    - `password` (string)

- `POST /api/v1/auth/login`
  - Required JSON body:
    - `email` (string)
    - `password` (string)

### Expenses

All expense endpoints require the header:

- `X-User-ID: <user-id>`

Use this header for all expense routes below.

- `POST /api/v1/expenses`
  - Create a new expense
  - Required JSON body:
    - `title` (string)
    - `amount` (number)
    - `category` (string)
    - `expense_date` (string, YYYY-MM-DD)
    - `note` (string)

- `GET /api/v1/expenses`
  - List expenses for the user
  - Optional filters:
    - `category`
    - `date_from` (YYYY-MM-DD)
    - `date_to` (YYYY-MM-DD)
    - `sort_by` (`amount` or `expense_date`)
    - `sort_order` (`asc` or `desc`)
    - `limit` (integer)

- `GET /api/v1/expenses/summary`
  - Get total spending summary for a date range
  - Required query params:
    - `date_from` (YYYY-MM-DD)
    - `date_to` (YYYY-MM-DD)

- `GET /api/v1/expenses/:id`
  - Get one expense by ID

- `PUT /api/v1/expenses/:id`
  - Update an expense by ID
  - Required JSON body same as create expense

- `DELETE /api/v1/expenses/:id`
  - Delete an expense by ID

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

### List expenses

```bash
curl -X GET "http://localhost:8080/api/v1/expenses?category=Food&date_from=2026-06-01&date_to=2026-06-30&sort_by=amount&sort_order=asc&limit=20" \
  -H "X-User-ID: 1"
```

### Get expense by ID

```bash
curl -X GET http://localhost:8080/api/v1/expenses/123 \
  -H "X-User-ID: 1"
```

### Update expense

```bash
curl -X PUT http://localhost:8080/api/v1/expenses/123 \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "title": "Dinner",
    "amount": 18.00,
    "category": "Food",
    "note": "Team dinner",
    "expense_date": "2026-06-02"
  }'
```

### Delete expense

```bash
curl -X DELETE http://localhost:8080/api/v1/expenses/123 \
  -H "X-User-ID: 1"
```

### Expense summary

```bash
curl -X GET "http://localhost:8080/api/v1/expenses/summary?date_from=2026-06-01&date_to=2026-06-30" \
  -H "X-User-ID: 1"
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
