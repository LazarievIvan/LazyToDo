# LazyToDo

Simple RESTful Todo List service in Go that allows users to manage todo items through a web API.
Built with **Go**, **PostgreSQL**, **Docker**, and **Swagger** documentation.

## Features
- Create, view, update, and delete to-do items.
- Query todos with optional filters: `status`, `orderBy`+`asc/desc`, `limit` and `page`.
- OpenAPI/Swagger support.
- Database schema migration with `golang-migrate`.
- Easy Docker setup.

## Tech Stack
- **Backend:** Go 1.24
- **Database:** PostgreSQL
- **ORM/DB Tooling:** `sqlc`
- **Migrations:** `golang-migrate`
- **API Docs:** Swagger (OpenAPI 3.0)

## Getting Started

### Prerequisites
- Docker + Docker Compose installed.

### Running the project

```bash
docker-compose up --build
```

---
## Testing
Perform testing using Postman or Swagger.

### Postman
Import postman collection, found in repository and refer to **API Endpoints**.

### Swagger
Swagger can be accessed at `http://localhost:8080/swagger/index.html`. From there it's straightforward.

---
## API Endpoints

| Method | Path             | Description                       |
|:-------|:------------------|:----------------------------------|
| POST   | `/add`            | Create a new todo item. Expects JSON body with `description` and `status`. |
| GET    | `/todos`          | Get all todos. Supports query params: `status`, `orderBy`, `asc`, `limit`, `page`. |
| GET    | `/todos/:id`      | Get a todo item by ID. |
| PUT    | `/todos/:id`      | Update a todo item by ID. JSON body can have `description` and/or `status`. |
| DELETE | `/todos/:id`      | Delete a todo item by ID. |

---
## Project Structure

```
├── cmd/
│   └── todo/                     # Main application: entry point and Swagger docs
│       ├── docs/                 # Swagger/OpenAPI
│       └── main.go               # Application startup (server initialization)
│
├── internal/
│   ├── db/
│   │   └── queries/               # SQL queries for sqlc code generation
│   │
│   ├── handler/
│   │   ├── routes.go              # HTTP routes setup (Gin router)
│   │   └── handler.go             # HTTP handlers for business logic
│   │
│   ├── models/
│   │   └── todo.go                # Structs representing application data (To-Dos)
│   │   └── params.go              # Structs representing query parameters (Sorting/Filtering/Pagination)
│   │
│   ├── repository/                # SQLC generated code and DB access layer
│   │   └── todos_repository.go    # DB access layer using sqlc generated and custom code
│   │
│   ├── server/
│       └── server.go              # HTTP server setup and configuration
│
├── migrations/                    # SQL migration files (for creating tables, etc.)
├── Dockerfile                     # Docker instructions for building the app container
├── docker-compose.yml             # Defines services (app + PostgreSQL) for easy startup
├── go.mod                         # Go modules: dependency list
├── go.sum                         # Go modules: dependency checksums
└── README.md                      # Project overview and setup guide
```

---
### Quick Navigation

| Folder         | Purpose                                           |
|:---------------|:--------------------------------------------------|
| `cmd/todo/`    | Where the app starts (main.go + swagger docs).    |
| `internal/handler/` | API endpoints and request routing.            |
| `internal/models/` | Data models used across the app.               |
| `internal/repository/` | Interactions with PostgreSQL database.    |
| `internal/db/` | SQL migrations and SQL queries for sqlc.           |
| `migrations/`  | Database migration files (used during startup).    |

---
