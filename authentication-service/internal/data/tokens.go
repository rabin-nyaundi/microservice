package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"log"
	"time"
)

// ScopeActivation indcates an activation token
// ScopeAuthentication indicates an authentication token
const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

// Token structure to hold data for 1 token from the database
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Scope     string    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

// TokenModel wraps the connection pool
type TokenModel struct {
	DB *sql.DB
}

//  New inserts a new token to the database
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := GenerateToken(userID, ttl, scope)

	if err != nil {
		log.Panic(err)
		return nil, err
	}
	err = m.Insert(token)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return token, nil
}

// Insert new token into the database
func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id`

	args := []interface{}{
		token.Hash,
		token.UserID,
		token.Expiry,
		token.Scope,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&token.UserID)
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

// DeleteAllForUser deletes all activation tokens for the user

// GenerateToken generates a new token
func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
