package model

import "golang.org/x/crypto/bcrypt"

// User
type User struct {
	ID                int
	Email             string
	Password          string
	EncryptedPassword string
	UserRole          int
	Active            bool
}

func (u *User) EncryptPasswordBeforeCreate() error {
	if len(u.Password) > 0 {
		encryptedString, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = encryptedString
	}

	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
