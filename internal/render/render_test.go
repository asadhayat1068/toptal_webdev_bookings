package render

import (
	"net/http"
	"testing"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	req, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(req.Context(), "flash", "123")
	result := AddDefaultData(&td, req)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	req, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww myWriter
	err = RenderTemplate(ww, req, "home", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing templates to the browser")
	}
	err = RenderTemplate(ww, req, "non-existed", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered non-existed template")
	}
}

func getSession() (*http.Request, error) {
	req, err := http.NewRequest("GET", "/uri", nil)
	if err != nil {
		return nil, err
	}
	ctx := req.Context()
	ctx, _ = session.Load(ctx, req.Header.Get("X-Session"))
	req = req.WithContext(ctx)
	return req, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error("Error writing templates to the browser")
	}
}
