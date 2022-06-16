package model

// User
type User struct {
	ID                int
	Email             string
	EncryptedPassword string
	UserRole          int
	Active            bool
}
