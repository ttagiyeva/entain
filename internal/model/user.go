package model

// UserDao is the domain object for users table.
type UserDao struct {
	ID      string  `db:"id"`
	Balance float32 `db:"balance"`
}
