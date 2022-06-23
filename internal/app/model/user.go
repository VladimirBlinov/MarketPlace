package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User
type User struct {
	ID                int
	Email             string
	Password          string
	EncryptedPassword string
	UserRole          int
	Active            bool
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(4, 50)),
	)
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
