package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is the structure that holds a user from the database
type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Active    bool      `json:"active"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Version   int       `json:"-"`
}

// password is the structure that hold a password
type password struct {
	plaintext *string
	hash      []byte
}

// UserModel structure that wraps the connection pool
type UserModel struct {
	DB *sql.DB
}

// GetAllUsers returns all users from the database
func (m UserModel) GetAllUsers() ([]*User, error) {
	return nil, nil
}

// GetByEmail returns a single user from the database by email
func (m UserModel) GetByEmail(email string) ([]*User, error) {
	return nil, nil
}

// GetOneUser returns a single user from the database by id
func (m UserModel) GetOneUser(id int) ([]*User, error) {
	return nil, nil
}

// Update updates and returns a single user
func (m UserModel) Update(user User) ([]*User, error) {
	return nil, nil
}

// Delete returns a single user from the database
func (m UserModel) Delete(user User) error {
	return nil
}

// Insert returns a single user inserted in to the database
func (m UserModel) Insert(user *User) ([]*User, error) {
	query := `
		INSERT INTO users (email, firstname, lastname, password_hash, active, role, version, created_at, updated_at)
		values($1, $2, $3, $4, false, 0, 0, $6, $7)
		RETURNING id`

	args := []interface{}{
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password.hash,
		user.Active,
		time.Now(),
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return nil, nil
}

// ResetPassword is a method called to change the user's password
func (m UserModel) ResetPassword(plaintext string, user *User) error {
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)

	if err != nil {
		log.Panic(err)
		return err
	}

	query := `
		UPDATE users 
		SET password_hash = $1
		WHER id = $2
		RETURNING id`

	args := []interface{}{
		newPasswordHash, user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)

	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

// Set method is called to encrypt user's passowrd
func (p *password) Set(plaintextPassword string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)

	if err != nil {
		log.Panic(err)
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// MatchesPassword method is called to check if plaintext passsword matches the hashed password
func (p *password) MatchesPassword(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
