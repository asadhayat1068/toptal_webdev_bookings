package main

import (
	"net/http"

	"github.com/asadhayat1068/toptal_webdev_bookings/pkg/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/pkg/handlers"
	"github.com/go-chi/chi/v4"
	"github.com/go-chi/chi/v4/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// // Using pat
	// mux := pat.New()
	// // Routes
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// Using chi
	mux := chi.NewRouter()
	// using Middlewares with chi
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	//Custom middleware
	mux.Use(WriteToConsole)
	// CSRF Middleware
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
