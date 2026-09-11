package user

import "database/sql"

type UserRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}
