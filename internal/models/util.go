package models

import "errors"

var (
	ErrNoRecord = errors.New("err: no matching record found")

	ErrInvalidCredentilas = errors.New("err: invalid credentials")

	ErrExsistingCrednetials = errors.New("err: dupliacte Credentials")
)
