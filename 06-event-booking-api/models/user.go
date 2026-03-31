package models

import (
	"time"

	"example.com/event-booking-api/auth"
	"example.com/event-booking-api/db"
	"example.com/event-booking-api/errs"
	"example.com/event-booking-api/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required,notblank,max=100"`
	CreatedAt time.Time `json:"createdAt"`
	UserCredentials
}

type UserCredentials struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password,omitempty" binding:"required,notblank,max=64"`
}

func (u *User) Save() error {
	u.CreatedAt = time.Now().UTC()
	utils.Trim(&u.Name, &u.Email)

	hashedPassword, err := auth.HashPassword(u.Password)

	if err != nil {
		return err
	}

	insert := `
	INSERT INTO users (name, created_at, email, password)
	VALUES (?, ?, ?, ?)
	`

	result, err := db.DB.Exec(insert, u.Name, u.CreatedAt, u.Email, hashedPassword)

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

	u.ID = id
	u.Password = ""

	return nil
}

func (u *User) Delete() error {
	delete := "DELETE FROM users WHERE id = ?"
	result, err := db.DB.Exec(delete, u.ID)

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return errs.ErrNotExists
	}

	return nil
}

func (uc *UserCredentials) GetAuthenticatedUser() (*User, error) {
	var u User
	utils.Trim(&uc.Email)

	query := `
	SELECT id, name, created_at, email, password
	FROM users
	WHERE email = ?
	`

	row := db.DB.QueryRow(query, uc.Email)
	err := row.Scan(&u.ID, &u.Name, &u.CreatedAt, &u.Email, &u.Password)

	if err != nil {
		if db.IsErrNoRows(err) {
			return nil, errs.ErrInvalidCredentials
		}

		return nil, err
	}

	if err := auth.ValidatePassword(u.Password, uc.Password); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	u.Password = ""
	return &u, nil
}

func GetUserByID(id int64) (*User, error) {
	var u User

	query := `
	SELECT id, name, created_at, email
	FROM users
	WHERE id = ?
	`

	row := db.DB.QueryRow(query, id)
	err := row.Scan(&u.ID, &u.Name, &u.CreatedAt, &u.Email)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
