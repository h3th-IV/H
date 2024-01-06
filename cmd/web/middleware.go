package main

import (
	"fmt"
	"net/http"
)

// middlerware to set security headers for all routes
func secureHaedersMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//(CSP)restricts location to load webPage resources
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		//used to strip off info that are sensitive in referer header when navigate to another URL
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		//do not sniff content type of response
		w.Header().Set("X-Content-Type-Options", "nosniff")
		//used  prevent clickjacking in old browsers
		w.Header().Set("X-Frame-Options", "deny")
		//set XSS peotection to disabled
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// middleware to log request
func (hb *hootBox) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hb.infolog.Printf("%v -%v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// middleware to handle errors incse of panic
func (hb *hootBox) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred function (will run in the event of a panic).
		defer func() {
			//use recover() to check if panic has occur
			if err := recover(); err != nil {
				//set header connection to close
				w.Header().Set("Connection", "close")
				//write 500 internal server error
				hb.serverErr(w, fmt.Errorf("%v", err))
			}
		}() //just like anormal defer line --> defer funcName()

		next.ServeHTTP(w, r)
	})
}
