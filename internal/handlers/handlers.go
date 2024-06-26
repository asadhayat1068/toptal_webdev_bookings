package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/config"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/driver"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/forms"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/helpers"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/render"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/repository"
	"github.com/asadhayat1068/toptal_webdev_bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v4"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (p *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home", &models.TemplateData{})
}

// About renders the about page
func (p *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about", &models.TemplateData{})
}

// Reservation renders the reservation page
func (p *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := p.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		p.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	room, err := p.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		p.App.Session.Put(r.Context(), "error", "cannot find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.Room.RoomName = room.RoomName
	reservation.Room.ID = reservation.RoomID

	p.App.Session.Put(r.Context(), "reservation", reservation)
	startDate := reservation.StartDate.Format("2006-01-02")
	endDate := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	data := make(map[string]any)
	data["reservation"] = reservation
	render.Template(w, r, "make-reservation", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (p *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation := models.Reservation{}
	// reservation, ok := p.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// if !ok {
	// 	log.Println("unable to load reservation data from session")
	// 	p.App.Session.Put(r.Context(), "error", "unable to load reservation data from session")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		p.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		log.Println("can't parse start date", err)
		p.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		log.Println("can't parse end date", err)
		p.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	form := forms.New(r.PostForm)

	form.Required("first_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})
		log.Println("Invalid Form data")
		return
	}

	reservationId, err := p.DB.InsertReservation(reservation)
	if err != nil {
		log.Println("can't insert reservation into database", err)
		p.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: reservationId,
		RestrictionID: 1,
	}

	_, err = p.DB.InsertRoomRestriction(restriction)
	if err != nil {
		log.Println("can't insert room restriction", err)
		p.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Send Notifications
	// 1. Send to guest
	htmlMessage := fmt.Sprintf(`
		<h3>Reservation Confirmation</h3>
		<p>Dear %s, </p>
		<p>This is confirm your reservation from %s to %s.</p>
	`,
		reservation.FirstName,
		reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"),
	)
	ml := models.MailData{
		To:      reservation.Email,
		From:    "owner123@bookings.pk",
		Subject: "Reservation confirmation!",
		Content: htmlMessage,
	}
	p.App.MailChan <- ml
	// 2. Send to owner

	p.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the generals page
func (p *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals", &models.TemplateData{})
}

// Majors renders the majors page
func (p *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors", &models.TemplateData{})
}

// Availability renders the availability page
func (p *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability", &models.TemplateData{})
}

// PostAvailability renders the availability page
func (p *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := p.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		p.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]any)
	data["rooms"] = rooms
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	p.App.Session.Put(r.Context(), "reservation", reservation)
	render.Template(w, r, "choose-room", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles requests for availability and send JSON response
func (p *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}

	available, _ := p.DB.SearchAvailabilityByDatesByRoomID(roomID, startDate, endDate)
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}
	respAsJson, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respAsJson)

}

// Contact renders the availability page
func (p *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact", &models.TemplateData{})
}

// ReservationSummary
func (p *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := p.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		p.App.ErrorLog.Println("Cannot get reservation data from session")
		log.Println("Cannot get reservation data from session")
		p.App.Session.Put(r.Context(), "error", "Cannot find reservation data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	p.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]any)
	data["reservation"] = reservation

	startDate := reservation.StartDate.Format("2006-01-02")
	endDate := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	render.Template(w, r, "reservation-summary", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom
func (p *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation, ok := p.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	reservation.RoomID = roomID
	p.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom
func (p *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("start")
	ed := r.URL.Query().Get("end")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}

	room, err := p.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}

	var reservation models.Reservation
	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.Room.RoomName = room.RoomName

	p.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (p *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLogin handles logging the user in
func (p *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = p.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if err != nil {
		log.Println(err)
		p.App.Session.Put(r.Context(), "error", "invalid input data")
		render.Template(w, r, "login", &models.TemplateData{
			Form: form,
		})
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if !form.Valid() {
		p.App.Session.Put(r.Context(), "error", "invalid input data")
		render.Template(w, r, "login", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := p.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		p.App.Session.Put(r.Context(), "error", "invalid login credentials")
		render.Template(w, r, "login", &models.TemplateData{
			Form: form,
		})
		return
	}
	p.App.Session.Put(r.Context(), "user_id", id)
	p.App.Session.Put(r.Context(), "success", "login successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs the user out
func (p *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	p.App.Session.Destroy(r.Context())
	p.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard shows the admin dashboard
func (p *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard", &models.TemplateData{})
}
