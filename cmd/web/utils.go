package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime/debug"
)

// write error err and stack trace  to the errlog attribute,
// sends 500 internal server error to user
// It returns a formatted stack trace of the current goroutine.
// A stack trace is a list of the function calls that have been
// executed up to a certain point in the program.
func (hb *hootBox) serverErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	hb.errlog.Output(2, trace) // track the debug stack 2 steps backward
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends stsatus code to user and corresponding status message
func (hb *hootBox) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 404-like responses  to user
func (hb *hootBox) notFoundErr(w http.ResponseWriter) {
	hb.clientErr(w, http.StatusNotFound)
}

func tests() error {
	Tra := sql.Tx{}
	Tra.Commit()
	return nil
}
