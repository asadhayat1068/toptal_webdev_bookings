package handlers

import (
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

func (p *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	p.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderPage(w, "home", &models.TemplateData{})
}

func (p *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello Again!"
	stringMap["remote_ip"] = p.App.Session.GetString(r.Context(), "remote_ip")
	render.RenderPage(w, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}
