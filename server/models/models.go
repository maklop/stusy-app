package models

type User struct {
	ID        uint
	Email     string
	Password  []byte
	Token     string
	ExpiresAt int64
}

type UserData struct {
	ID        uint
	UserID    uint
	FirstName string
	LastName  string
}
