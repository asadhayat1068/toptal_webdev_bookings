package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	log.Println("Server is running at port", PORT)

	srv := http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() (*driver.DB, error) {
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
	//Connect to database
	log.Println("Connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=asad password=")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Successfully connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.UseCache = false
	app.TemplateCache = tc

	// Init Handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// Init render
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	// Routes
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	return db, nil
}
