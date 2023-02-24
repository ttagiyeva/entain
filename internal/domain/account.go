package domain

// User is the domain object for users table.
type User struct {
	ID      string  `db:"id"`
	Balance float32 `db:"balance"`
}
