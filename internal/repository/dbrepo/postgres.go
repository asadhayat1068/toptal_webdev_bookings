package dbrepo

import (
	"context"
	"time"

	"github.com/asadhayat1068/toptal_webdev_bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int
	stmt := `insert into reservations (
		first_name,
		last_name,
		email,
		phone,
		start_date,
		end_date,
		room_id,
		created_at,
		updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, err
}

// InsertRoomRestriction inserts a room restriction record in database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int
	stmt := `insert into room_restrictions (
		start_date,
		end_date,
		room_id,
		reservation_id,
		created_at,
		updated_at,
		restriction_id) values (
			$1,$2,$3,$4,$5,$6,$7
		) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.CreatedAt,
		r.UpdatedAt,
		r.RestrictionID,
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}
