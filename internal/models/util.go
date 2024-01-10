package models

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrNoRecord = errors.New("err: no matching record found")

	ErrInvalidCredentilas = errors.New("err: invalid credentials")

	ErrExsistingCrednetials = errors.New("err: dupliacte Credentials")
	MySQLErr                *mysql.MySQLError
)
