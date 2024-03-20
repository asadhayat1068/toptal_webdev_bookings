package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/driver"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func getRoutes() http.Handler {
	// What to store in Session
	gob.Register(models.Reservation{})
	// Init App Configs
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Persist = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	InfoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	ErrorLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = InfoLog
	app.ErrorLog = ErrorLog
	app.TemplateCache = tc
	app.UseCache = true
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=asad password=")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Successfully connected to database")
	// Init Handlers
	repo := NewRepo(&app, db)
	NewHandlers(repo)

	// Init render
	render.NewRenderer(&app)
	mux := chi.NewRouter()
	// using Middlewares with chi
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	//Custom middleware
	mux.Use(WriteToConsole)
	// CSRF Middleware
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suits", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.Contact)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

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

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files ending with .page.tmpl from ./templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	// range through all pages
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}

	return myCache, nil
}
