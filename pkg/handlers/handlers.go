package handlers

import (
	"fmt"
	"net/http"

	"github.com/asadhayat1068/toptal_webdev_bookings/pkg/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/pkg/models"
	"github.com/asadhayat1068/toptal_webdev_bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (p *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	p.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home", &models.TemplateData{})
}

// About renders the about page
func (p *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello Again!"
	stringMap["remote_ip"] = p.App.Session.GetString(r.Context(), "remote_ip")
	render.RenderTemplate(w, r, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the reservation page
func (p *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{})
}

// Generals renders the generals page
func (p *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals", &models.TemplateData{})
}

// Majors renders the majors page
func (p *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors", &models.TemplateData{})
}

// Availability renders the availability page
func (p *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability", &models.TemplateData{})
}

// PostAvailability renders the availability page
func (p *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

// Contact renders the availability page
func (p *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact", &models.TemplateData{})
}
