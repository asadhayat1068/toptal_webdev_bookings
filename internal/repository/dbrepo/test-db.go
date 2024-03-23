package dbrepo

import (
	"errors"
	"time"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction record in database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) (int, error) {

	return 1, nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for room and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(roomID int, start, end time.Time) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of all available rooms for a date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID Gets a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("No room found")
	}

	return room, nil
}
