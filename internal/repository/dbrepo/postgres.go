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

// SearchAvailabilityByDatesByRoomID returns true if availability exists for room and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(roomID int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `select 
		count(id) 
		from room_restrictions 
		where 
			room_id = $1
			and $2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of all available rooms for a date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `select 
		r.id, r.room_name
		from rooms r 
		where 
			r.id not in (
				select rr.room_id from room_restrictions rr where
				$1 < rr.end_date and $2 > rr.start_date
			);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, nil
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetRoomByID Gets a room by ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}
