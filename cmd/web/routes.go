package main

import (
	"net/http"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/handlers"
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
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suits", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
