package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// What to store in Session
	gob.Register(models.Reservation{})
	// Init App Configs
	testApp.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	testApp.Session = session
	app = &testApp
	os.Exit(m.Run())
}

type myWriter struct{}

func (w myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (w myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (w myWriter) WriteHeader(statusCode int) {
}
