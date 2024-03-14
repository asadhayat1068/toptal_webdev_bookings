package main

import (
	"testing"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/go-chi/chi/v4"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("Expected Type *chi.Mux, Got %T", v)

	}
}
