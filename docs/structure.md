# Project structure


#### **Root**
README.md  
git-related  
go.sum  
go.mod  
Docker files  

#### **Docs**  
#### **Migrations**  
#### **Web**  

#### **Cmd**  
main.go  

#### **Internal**  
**config**  
└── config.go  

**database**  
├── database.go  

**handlers**  
├── init.go  

**entities**  
├── category  
│   ├── handlers.go  
│   ├── model.go  
│   ├── repository.go  
│   └── service.go  
├── thread  
│   ├── handlers.go  
│   ├── model.go  
│   ├── repository.go  
│   └── service.go  
└── user  
│   ├── handlers.go  
│   ├── model.go  
│   ├── repository.go  
│   └── service.go  

**routes**  
└── routes.go  

**middleware**  
├── logging.go  
└── recovery.go  




