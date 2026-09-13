# Task Manager API

A RESTful task management API built with Go and PostgreSQL.

The project was built from scratch using Go's standard `net/http` package, with a layered architecture separating HTTP handlers, models, database access, and repositories.

## Features

- Create, read, update, and delete tasks
- Partial task updates using `PATCH`
- PostgreSQL persistence
- JSON request/response handling
- Request validation
- HTTP status code handling
- Unknown JSON field rejection
- Repository pattern for database access
- Unit tests
- PostgreSQL integration tests
- HTTP handler tests
- Bruno API collection for manual API testing

## Tech Stack

- **Go**
- **PostgreSQL**
- **net/http**
- **encoding/json**
- **pgx/v5**
- **Bruno**
- **Go testing package**

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | Get all tasks |
| POST | `/tasks` | Create a task |
| GET | `/tasks/:id` | Get a task by ID |
| PUT | `/tasks/:id` | Replace a task |
| PATCH | `/tasks/:id` | Partially update a task |
| DELETE | `/tasks/:id` | Delete a task |

## Project Structure

```text
taskmanager/
├── main.go
├── model/
├── repository/
├── database/
├── migrations/
│   └── 001_create_tasks.sql
├── bruno/
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```

## Database

The API uses PostgreSQL for persistent storage.

The `tasks` table contains:

- `id`
- `title`
- `description`
- `completed`
- `created_at`
- `updated_at`

Database migrations are stored in the `migrations/` directory.

## Running Locally

### Prerequisites

- Go
- PostgreSQL

### 1. Clone the repository

```bash
git clone <repository-url>
cd taskmanager
```

### 2. Configure environment variables

Create a `.env` file based on `.env.example`:

```env
DATABASE_URL=postgres://taskmanager_user:password@127.0.0.1:5432/taskmanager
```

### 3. Set up the database

Create the PostgreSQL database and run the migration:

```bash
psql -U taskmanager_user -d taskmanager -f migrations/001_create_tasks.sql
```

### 4. Run the API

```bash
go run .
```

The API will start on:

```text
http://localhost:8080
```

## Testing

The project contains tests at multiple levels.

### Unit Tests

Tests validation, path handling, JSON decoding, and other application logic.

### Repository Integration Tests

Tests repository operations against a real PostgreSQL test database.

### HTTP Handler Tests

Tests API handlers independently using fake repositories, covering successful requests, validation errors, not-found responses, and repository failures.

Run the complete test suite with:

```bash
go test ./...
```

## API Testing with Bruno

A Bruno collection is included in the `bruno/` directory for manually testing the API endpoints.

The collection covers the complete CRUD workflow.

## Architecture

The project follows a simple layered structure:

```text
HTTP Request
     │
     ▼
  Handler
     │
     ▼
   Model
     │
     ▼
 Repository
     │
     ▼
 PostgreSQL
```

The repository layer handles database operations, while handlers are responsible for HTTP-specific concerns such as request decoding, validation, response formatting, and status codes.

Small interfaces are used by handlers so they can be tested independently from the PostgreSQL implementation.

## What I Learned

- Building REST APIs with Go's standard library
- HTTP routing and request handling
- JSON encoding and decoding
- PostgreSQL integration using `pgx`
- Repository-based database access
- Dependency injection through interfaces
- Unit and integration testing
- HTTP handler testing with fake dependencies
- API testing with Bruno
- Designing and organizing a Go project without a web framework
