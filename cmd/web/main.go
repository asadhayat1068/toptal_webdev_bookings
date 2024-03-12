package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/handlers"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/render"
)

const PORT = ":8080"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	// Init App Configs
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Persist = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.UseCache = false
	app.TemplateCache = tc

	// Init Handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// Init render
	render.NewTemplates(&app)

	// Routes
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Server is running at port", PORT)
	// http.ListenAndServe(PORT, nil)

	srv := http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
