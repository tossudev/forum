# Architecture

## Overview 
- [App](#app)
- [Routes](#routes)
- [Middleware](#middleware)
- [Entities](#entities)
    - [Entity Handler](#entity-handler)
    - [Entity Service](#entity-service)
    - [Entity Repository](#entity-repository) 

### Data flow diagram
       HTTP Request
            │
            ▼
          Routes
            │
            ▼
        Middleware
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

## App
[Back to top](#overview)
```text
internal/routes/init_handlers.go
```
#### Responsibilities
- Create handlers for entities
- Wrap handlers in middleware (found in internal/middleware)

### Example

```golang
type App struct {
	User    *user.Handler
	Thread  *thread.Handler, etc...
}

func New(db *sql.DB) http.Handler {
    userRepo := user.NewRepo(db)
    userService := user.NewService(userRepo)
    userHandler := user.NewHandler(userService)

    threadRepo := thread.NewRepository(db)
    threadService := thread.NewService(threadRepo)
    threadHandler := thread.NewHandler(threadService)
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
[Back to top](#overview)
```text
internal/routes/routes.go
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
[Back to top](#overview)
```text
internal/middleware/
//except where noted
```

Middleware intercepts every incoming HTTP request before it reaches the handlers.

- **Timeout** - Creates a 5-second timeout for every request.
- **User Middleware** <em>(In the user entity)</em> - Gets user information and attaches it to the context. 
- **Session Authentication** <em>(In the session package)</em> - Checks cookie for session id and validates it. 
- **Time Format** <em>(In the timeformat package)</em> - Gets user timezone. 
- **Request Logging** – Logs important information about the requests.
- **Panic Recovery** – Recovers from unexpected panics and returns a standardized `500 Internal Server Error` response instead of terminating the application.

---

## Entities
[Back to top](#overview)
```text
internal/entities/
```
Entities contain layers for each database entity  

### Entities Used
- Category
- Comment
- Image
- Like <em>(also handles dislikes)</em>
- Thread
- User

### Example Entity Layer
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
userHandler.Create() //handles http requests & responses
        │
        ▼
userService.Create() //handles validation & formatting
        │
        ▼
userRepo.Create() //handles database logic
```

---

## Entity Handler
[Back to top](#overview)
```text
internal/entities/<entity>/handler.go
```

### Responsibilities

Handlers are responsible for processing HTTP requests and responses.

Responsibilities include:

- Parsing JSON request bodies
- Reading path parameters
- Reading query parameters
- Calling the appropriate service
- Returning JSON responses
- Setting HTTP status codes

---

## Entity Service
[Back to top](#overview)
```text
internal/entities/<entity>/service.go
```

### Responsibilities

The service layer contains all business logic.

Examples include:

- Creating and updating entities
- Validation beyond simple field checks
- HTML Parsing
- Pagination

The service layer coordinates one or more repositories to complete an operation.

---

## Entity Repository
[Back to top](#overview)
```text
internal/entities/<entity>/repository.go
```

#### Responsibilities

The repository layer interacts directly with the SQLite database.

Responsibilities include:

- CRUD operations
- Raw SQL queries
- Join queries
- Returning repository-specific errors

### Example

```sql
SELECT *
FROM users
WHERE id = ?;
```

---

