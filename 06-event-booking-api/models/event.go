package models

import (
	"database/sql"
	"strings"
	"time"

	"example.com/event-booking-api/db"
	"example.com/event-booking-api/errs"
	"example.com/event-booking-api/utils"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required,notblank,max=100"`
	Description string    `json:"description" binding:"required,notblank,max=1000"`
	Location    string    `json:"location" binding:"required,notblank,max=250"`
	DateTime    time.Time `json:"dateTime" binding:"required"`
	UserID      int64     `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EventDetails struct {
	Event
	UserName            string `json:"userName"`
	LoggedUserHasBooked bool   `json:"loggedUserHasBooked"`
	LoggedUserBookingID *int64 `json:"loggedUserBookingId"`
}

func (e *Event) Save() error {
	e.DateTime = e.DateTime.UTC()
	e.CreatedAt = time.Now().UTC()
	e.UpdatedAt = time.Now().UTC()
	utils.Trim(&e.Name, &e.Description, &e.Location)

	insert := `
	INSERT INTO events (
		name,
		description,
		location,
		date_time,
		user_id,
		created_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.DB.Exec(
		insert,
		e.Name,
		e.Description,
		e.Location,
		e.DateTime,
		e.UserID,
		e.CreatedAt,
		e.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	e.ID = id
	return nil
}

func (e *Event) Update() error {
	e.DateTime = e.DateTime.UTC()
	e.UpdatedAt = time.Now().UTC()
	utils.Trim(&e.Name, &e.Description, &e.Location)

	update := `
	UPDATE events
	SET name = ?, description = ?, location = ?, date_time = ?, updated_at = ?
	WHERE id = ?
	`

	_, err := db.DB.Exec(
		update,
		e.Name,
		e.Description,
		e.Location,
		e.DateTime,
		e.UpdatedAt,
		e.ID,
	)

	return err
}

func (e *Event) Delete() error {
	delete := "DELETE FROM events WHERE id = ?"
	_, err := db.DB.Exec(delete, e.ID)

	return err
}

func GetAllEvents(lastID, hostID, guestID int64) ([]Event, error) {
	events := []Event{}

	var sb strings.Builder
	args := []interface{}{}

	sb.WriteString("SELECT ")
	sb.WriteString("	e.id, ")
	sb.WriteString("	e.name, ")
	sb.WriteString("	e.description, ")
	sb.WriteString("	e.location, ")
	sb.WriteString("	e.date_time, ")
	sb.WriteString("	e.user_id, ")
	sb.WriteString("	e.created_at, ")
	sb.WriteString("	e.updated_at ")
	sb.WriteString("FROM events e ")

	if guestID > 0 {
		sb.WriteString("JOIN bookings b ON e.id = b.event_id AND b.user_id = ? ")
		args = append(args, guestID)
	}

	sb.WriteString("WHERE 1 = 1 ")

	if lastID > 0 {
		sb.WriteString("AND e.id < ? ")
		args = append(args, lastID)
	}

	if hostID > 0 {
		sb.WriteString("AND e.user_id = ? ")
		args = append(args, hostID)
	}

	sb.WriteString("ORDER BY e.id DESC ")
	sb.WriteString("LIMIT 50")

	query := sb.String()
	rows, err := db.DB.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var e Event

		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Description,
			&e.Location,
			&e.DateTime,
			&e.UserID,
			&e.CreatedAt,
			&e.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventDetailsByID(id, loggedUserID int64) (*EventDetails, error) {
	var e EventDetails
	var bookingID sql.NullInt64

	query := `
	SELECT
		e.id,
		e.name,
		e.description,
		e.location,
		e.date_time,
		e.user_id,
		e.created_at,
		e.updated_at,
		u.name,
		b.id AS booking_id
	FROM events e
	JOIN users u ON e.user_id = u.id
	LEFT JOIN bookings b ON e.id = b.event_id AND b.user_id = ?
	WHERE e.id = ?
	`

	row := db.DB.QueryRow(query, loggedUserID, id)

	err := row.Scan(
		&e.ID,
		&e.Name,
		&e.Description,
		&e.Location,
		&e.DateTime,
		&e.UserID,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.UserName,
		&bookingID,
	)

	if err != nil {
		if db.IsErrNoRows(err) {
			return nil, errs.ErrNotExists
		}

		return nil, err
	}

	if bookingID.Valid {
		e.LoggedUserHasBooked = true
		e.LoggedUserBookingID = &bookingID.Int64
	} else {
		e.LoggedUserHasBooked = false
		e.LoggedUserBookingID = nil
	}

	return &e, nil
}
