package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
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

// middleware to handle errors incase of panic
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

// middleware to prevents access to routes that require auth
func (hb *hootBox) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//cehck if user session is logedIn
		if !hb.isAuthenticate(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the cache
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// func to mitigate CSRF request
func noCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

// middleware to
func (hb *hootBox) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//retrieve the authenticatedUserID from the session
		id := hb.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		//check if the user Exist in the database with the provided ID
		exists, err := hb.users.Exists(id)
		if err != nil {
			hb.serverErr(w, err)
			return
		}
		//if user exists, this means request is coming from users that exists in the DB
		//make a copy of the request(with isAuthenticatedContextKey set to true --
		//which shows user is uthenticated)
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedcontextKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
