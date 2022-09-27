package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	DuplicateEmail = errors.New("duplicate email found")
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
	query := `
		SELECT firtsname, lastname, email, active, role, version, created_at, updated_at
		FROM users ORDER BY lastname`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Active,
			&user.Role,
			&user.Version,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			log.Panic(err)
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		log.Panic(err)
		return nil, err
	}

	return users, nil
}

// GetByEmail returns a single user from the database by email
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, firstname, lastname, email, password_hash, active, role, version
		FROM users
		WHERE email = $1`

	var user User

	args := []interface{}{email}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.hash,
		&user.Active,
		&user.Role,
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `sql: no rows in result set`:
			return nil, ErrorRecordNotFound
		default:
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &user, nil
}

// GetOneUser returns a single user from the database by id
func (m UserModel) GetOneUser(id int) (*User, error) {
	query := `
		SELECT firstname, lastname, email, active, role, version
		FROM users
		WHERE id = $1`

	var user User

	args := []interface{}{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Active,
		&user.Role,
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `sql: no rows in result set`:
			return nil, ErrorRecordNotFound
		default:
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return &user, nil
}

// Update updates and returns a single user
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users SET 
		email = $1, 
		fristname=$2, 
		lastname=$3, 
		version=$4, 
		updated_at=$5
		WHERE id=$6`

	args := []interface{}{
		user.Email,
		user.FirstName,
		user.LastName,
		user.Version + 1,
		time.Now(),
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)

	if err != nil {
		return err
	}
	return nil
}

// Delete returns a single user from the database
func (m UserModel) Delete(user User) error {
	query := `
	DELETE * FROM users
	WHERE id = $1`

	args := []interface{}{
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

// Insert returns a single user inserted in to the database
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (email, firstname, lastname, password_hash, active, role, version, created_at, updated_at)
		values($1, $2, $3, $4, false, 0, 0, $5, $6)
		RETURNING id`

	args := []interface{}{
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password.hash,
		time.Now(),
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return DuplicateEmail
		default:
			log.Panic(err)
		}
		log.Panic(err)
		return err
	}
	return nil
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
