package db

import (
	"database/sql"
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "api.db?_pragma=foreign_keys(1)")

	if err != nil {
		panic("Could not connect to database.")
	}

	// Best settings for SQLite
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(0)

	createTables()
}

func IsErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsErrUniqueConstraint(err error) bool {
	var sqliteErr *sqlite.Error
	errCode := sqlite3.SQLITE_CONSTRAINT_UNIQUE

	return errors.As(err, &sqliteErr) && sqliteErr.Code() == errCode
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	`

	if _, err := DB.Exec(createUsersTable); err != nil {
		panic("Could not create users table: " + err.Error())
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		date_time DATETIME NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_events_user_id
	ON events(user_id);
	`

	if _, err := DB.Exec(createEventsTable); err != nil {
		panic("Could not create events table: " + err.Error())
	}

	createBookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_bookings_event_user
	ON bookings(event_id, user_id);

	CREATE INDEX IF NOT EXISTS idx_bookings_user_id
	ON bookings(user_id);
	`

	if _, err := DB.Exec(createBookingsTable); err != nil {
		panic("Could not create bookings table: " + err.Error())
	}
}
