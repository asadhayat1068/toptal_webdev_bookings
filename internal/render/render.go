package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
	"github.com/justinas/nosurf"
)

var pathToTemplates = "./templates"

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// create a template cache
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tmpl+".page.tmpl"]

	if !ok {
		return errors.New("Could not get template from cache")
	}

	buf := new(bytes.Buffer)
	AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}
	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser:", err)
		return err
	}

	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {
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
