# Architecture

## Overview

### Data flow diagram
       HTTP Request
            │
            ▼
          Routes
            │
            ▼
        Middleware
        (logging)
        (recovery)
            │
            ▼
      Entity handler
            │
            ▼
      Entity service
            │
            ▼
    Entity repository
            │
            ▼
         Database

---

# Layers

## Handlers
```text
internal/handlers/
```
#### Responsibilities
- Create handlers for entites
- Wrap handers in middleware

### Example

```golang
type App struct {
	User    *user.Handler
	Thread  *thread.Handler, etc...
}

func New(db *sql.DB) http.Handler {
    userRepo := user.NewRepo(db)
    userService := user.NewService(userRepo)
    userHandler := user.NewHandler(userService, validator.Validate)

    threadRepo := thread.NewRepository(db)
    threadService := thread.NewService(threadRepo)
    threadHandler := thread.NewHandler(threadService, validate)
    // etc...

    app := App{
        User:   userHandler,
        Thread: threadHandler,
    }

    mux := routes.NewRoutes(&app)

    handler := middleware.Logger(mux)
    handler = middleware.Recovery(handler)

    return handler
```

---

## Routes
```text
internal/routes/
```
#### Responsibilities
- Register all API endpoints
- Map HTTP methods and URLs to handlers

### Example

```golang
func NewRoutes(h Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", app.UserHandler.Get)
	mux.HandleFunc("POST /users", app.UserHandler.Create)

	return mux
}
```

---

## Middleware
```text
internal/middleware/
```

Middleware intercepts every incoming HTTP request before it reaches the handlers.

- **Request Logging** – Logs important information about the requests.
- **Panic Recovery** – Recovers from unexpected panics and returns a standardized `500 Internal Server Error` response instead of terminating the application.

---

## Entities
```text
internal/entities/
```

Entities contain layers for each database entity  
### Example entity layer
```bash
internal/entities/<entity>
└── model.go
└── handlers.go
└── service.go
└── repository.go
```

### Entity Example Flow

```text
Request & routing
        │
        ▼
userHandler.Create()
        │
        ▼
userService.Create()
        │
        ▼
userRepo.Create()
```

---

## Entity handler
```text
internal/entities/<entity>/handler.go
```

#### Responsibilities

Handlers are responsible for processing HTTP requests and responses.

Responsibilities include:

- Parsing JSON request bodies
- Reading path parameters
- Reading query parameters
- Calling the appropriate service
- Returning JSON responses
- Setting HTTP status codes

---

## Entity service
```text
internal/entities/<entity>/service.go
```

### Responsibilities

The service layer contains all business logic.

Examples include:

- Creating and updating entities
- Filtering
- Searching
- Pagination
- Validation beyond simple field checks

The service layer coordinates one or more repositories to complete an operation.

---

## Repository Layer
```text
internal/entities/<entity>/repository.go
```

#### Responsibilities

The repository layer interacts directly with the SQLite database.

Responsibilities include:

- CRUD operations
- Raw SQL queries
- Join queries
- Transactions
- Returning repository-specific errors

### Example

```sql
SELECT *
FROM users
WHERE id = ?;
```

---

