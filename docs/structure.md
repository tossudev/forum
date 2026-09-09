# Project structure

Keep this up to date
  
/forum  
├── cmd  
│   └── main.go  
├── internal  
│   ├── config  
│   │   └── config.go  
│   ├── database  
│   │   ├── database.go  
│   │   └── migrations  
│   │       └── 001_init.sql  
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
│   ├── http  
│   │   └── routes.go  
│   └── middleware  
│       ├── logging.go  
│       └── recovery.go  
└── web  
    └── index.html  
