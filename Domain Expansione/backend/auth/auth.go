package auth

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
	Password string
}

// HashPassword hashes the plain-text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares plain password with hashed one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RegisterUser inserts a new user into the database with a hashed password
func RegisterUser(db *sql.DB, email, password, username, phone string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (email, password, username, phone_number)
		VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(query, email, hashedPassword, username, phone)
	return err
}

// AuthenticateUser checks if email/password is correct and returns user ID and username
func AuthenticateUser(db *sql.DB, email, password string) (int, string, error) {
	var id int
	var username string
	var hashedPassword string

	query := `SELECT id, password, username FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&id, &hashedPassword, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New("user not found")
		}
		return 0, "", err
	}

	if !CheckPasswordHash(password, hashedPassword) {
		return 0, "", errors.New("invalid password")
	}

	return id, username, nil
}
