# Internal Transfer System

A secure internal transfers system written in Go, using PostgreSQL. This service enables financial transfers between accounts with:

---

## Features

-   Create and query accounts
-   Transfer money between accounts safely

---

## Requirements

-   Go 1.20+
-   PostgreSQL 13+
-   (Optional) Postman for API testing

---

## Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/your-org/internal-transfer-system.git
cd internal-transfer-system
```

### 2. Create database schema

```bash
psql -U your_user -d your_db -f db/schema.sql
```

### 3. Set the environment variable

```bash
export POSTGRES_DSN="postgres://username:password@localhost:5432/dbname?sslmode=disable"
```

### 4. Run the app

```bash
go run main.go
```

---

## Running Tests

Run all tests:

```bash
go test ./... -v
```

---

## Postman Collection

The `postman/` folder contains a Postman collection for testing all endpoints.

**Steps**:

1. Open Postman
2. Import: `postman/internal-transfer-system.postman_collection.json`
3. Test each endpoint (create, get, transfer)
