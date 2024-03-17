package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/handlers"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/helpers"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/render"
)

const PORT = ":8080"

var app config.AppConfig
var InfoLog *log.Logger
var ErrorLog *log.Logger
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is running at port", PORT)
	// http.ListenAndServe(PORT, nil)

	srv := http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {
	// What to store in Session
	gob.Register(models.Reservation{})
	// Init App Configs
	app.InProduction = false

	InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = InfoLog
	app.ErrorLog = ErrorLog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Persist = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return err

	}
	app.UseCache = false
	app.TemplateCache = tc

	// Init Handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// Init render
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	// Routes
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	return nil
}
