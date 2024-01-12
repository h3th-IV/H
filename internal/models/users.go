package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users(name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP()) `
	//create dB transacions
	trx, err := um.DB.Begin()
	if err != nil {
		return err
	}
	defer trx.Rollback()

	//statement
	stmt, err := trx.Prepare(query)
	if err != nil {
		return err
	}

	//execute statement
	_, err = stmt.Exec(name, email, string(hash))
	if err != nil {
		//check if err is a mysql error type
		if errors.As(err, &MySQLErr) {
			//check if error is existing credentials (not unique) with the constraint 'users_uc_email'
			if MySQLErr.Number == 1062 && strings.Contains(MySQLErr.Message, "users_uc_email") {
				return ErrExsistingCrednetials
			}
		}
		return err
	}
	trx.Commit()
	return nil

}

// check if a user exist in the DB
func (um *UserModels) Exists(id int) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`
	err := um.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}

// authenticate user if the if email match password
func (um *UserModels) Authenticate(email, password string) (int, error) {
	//retrieve user from DB if their email exist
	query := `SELECT id, hashed_password FROM users WHERE email = ?`
	tx, err := um.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	user := User{}
	row := stmt.QueryRow(email)
	err = row.Scan(&user.ID, &user.HashedPassowrd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	//check whether plain text password matches retrieved hash
	err = bcrypt.CompareHashAndPassword(user.HashedPassowrd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (um *UserModels) CheckEmail(email string) bool {
	query := `SELECT * FROM users WHERE email = ?`

	//use transaction pool
	tx, err := um.DB.Begin()
	if err != nil {
		return false
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return false
	}

	row := stmt.QueryRow(email)
	tx.Commit()
	return row != nil
}
