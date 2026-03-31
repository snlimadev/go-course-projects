package models

import (
	"time"

	"example.com/event-booking-api/db"
	"example.com/event-booking-api/errs"
)

type Booking struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"eventId"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type BookingDetails struct {
	Booking
	UserName string `json:"userName"`
}

func (b *Booking) Save() error {
	b.CreatedAt = time.Now().UTC()

	insert := `
	INSERT INTO bookings (event_id, user_id, created_at)
	VALUES (?, ?, ?)
	`

	result, err := db.DB.Exec(insert, b.EventID, b.UserID, b.CreatedAt)

	if err != nil {
		if db.IsErrUniqueConstraint(err) {
			return errs.ErrAlreadyExists
		}

		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	b.ID = id
	return nil
}

func (b *Booking) Delete() error {
	delete := "DELETE FROM bookings WHERE id = ?"
	_, err := db.DB.Exec(delete, b.ID)

	return err
}

func GetBookingsCountByEventID(eventID int64) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM bookings WHERE event_id = ?"

	if err := db.DB.QueryRow(query, eventID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func GetAllBookingsByEventID(eventID, lastID int64) ([]BookingDetails, error) {
	bookings := []BookingDetails{}

	query := `
	SELECT b.id, b.event_id, b.user_id, b.created_at, u.name
	FROM bookings b
	JOIN users u ON b.user_id = u.id
	WHERE b.event_id = ?
	AND b.id > ?
	ORDER BY b.id
	LIMIT 50
	`

	rows, err := db.DB.Query(query, eventID, lastID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var b BookingDetails
		err := rows.Scan(&b.ID, &b.EventID, &b.UserID, &b.CreatedAt, &b.UserName)

		if err != nil {
			return nil, err
		}

		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func GetBookingByID(id int64) (*Booking, error) {
	var b Booking

	query := `
	SELECT id, event_id, user_id, created_at
	FROM bookings
	WHERE id = ?
	`

	row := db.DB.QueryRow(query, id)
	err := row.Scan(&b.ID, &b.EventID, &b.UserID, &b.CreatedAt)

	if err != nil {
		if db.IsErrNoRows(err) {
			return nil, errs.ErrNotExists
		}

		return nil, err
	}

	return &b, nil
}
