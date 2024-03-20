package repository

import "github.com/asadhayat1068/toptal_webdev_bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) (int, error)
}
