# Project structure

Keep this up to date

 ```sh
Cmd
│   └── main.go  
│  
Internal  
│   ├── config  
│   │   └── config.go  
│   │  
│   ├── database  
│   │   ├── database.go  
│   │   └── migrations  
│   │       └── 001_init.sql (schema)  
│   │  
│   ├── entities  
│   │   ├── category  
│   │   │   ├── handlers.go  
│   │   │   ├── model.go  
│   │   │   ├── repository.go  
│   │   │   └── service.go  
│   │   ├── thread  
│   │   │   ├── handlers.go  
│   │   │   ├── model.go  
│   │   │   ├── repository.go  
│   │   │   └── service.go  
│   │   └── user  
│   │       ├── handlers.go  
│   │       ├── model.go  
│   │       ├── repository.go  
│   │       └── service.go  
│   │  
│   ├── http  
│   │   └── routes.go  
│   └── middleware  
│       ├── logging.go  
│       └── recovery.go  
│  
Web  
    └── index.html
```
