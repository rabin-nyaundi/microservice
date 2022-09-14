package data

import "database/sql"


// Models structs that wraps the models
type Models struct {
	User UserModel
}

// NewModel returns models struct with initialized models
func NewModel(db *sql.DB) Models {
	return Models{
		User: UserModel{DB: db},
	}
}
