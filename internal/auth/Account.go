package auth

type Account struct {
	ID        int64  `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
}
