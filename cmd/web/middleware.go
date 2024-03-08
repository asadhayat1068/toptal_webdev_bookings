package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// Custom Middleware
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I am logging from custom middleware")
		next.ServeHTTP(w, r)
	})
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad load and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	_next := session.LoadAndSave(next)
	return _next
}
