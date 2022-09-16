package data

import (
	"database/sql"
	"errors"
)

// ErrorRecordNotFound returns record not found error
var (
	ErrorRecordNotFound = errors.New("record not found")
)

// Models structs that wraps the models
type Models struct {
	User  UserModel
	Token TokenModel
}

// NewModel returns models struct with initialized models
func NewModel(db *sql.DB) Models {
	return Models{
		User:  UserModel{DB: db},
		Token: TokenModel{DB: db},
	}
}
