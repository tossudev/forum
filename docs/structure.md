 # Project structure

  Keep this up to date

  ```sh
  cmd
  │   └── main.go
  │
  internal
  │   ├── config
  │   │   └── config.go
  │   │
  │   ├── database
  │   │   └── database.go
  │   │
  │   ├── entities
  │   │   ├── category
  │   │   │   ├── handlers.go
  │   │   │   ├── model.go
  │   │   │   ├── repository.go
  │   │   │   └── service.go
  │   │   ├── comment
  │   │   │   ├── handlers.go
  │   │   │   ├── model.go
  │   │   │   ├── repository.go
  │   │   │   └── service.go
  │   │   ├── thread
  │   │   │   ├── handlers.go
  │   │   │   ├── model.go
  │   │   │   ├── repository.go
  │   │   │   └── service.go
  │   │   └── user
  │   │       ├── handlers.go
  │   │       ├── model.go
  │   │       ├── repository.go
  │   │       └── service.go
  │   │
  │   ├── http
  │   │   └── routes.go
  │   └── middleware
  │       ├── logging.go
  │       └── recovery.go
  │
  migrations
  │   └── 001_init.sql (schema)
  │
  web
      └── index.html
  ```

