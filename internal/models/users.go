package models

import (
	"database/sql"
	"time"
)

// user modell
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassowrd []byte
	CreateAt       time.Time
}

// define UserModel which wraps the DB connection pool
type UserModels struct {
	DB *sql.DB
}

// create new user here
func (um *UserModels) Insert(name, email, password string) error {
	return nil
}

// authenticate user if the if email match password
func (um *UserModels) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// to check is a user with particular ID exist
func (um *UserModels) Exists(id int) (bool, error) {
	return false, nil
}
