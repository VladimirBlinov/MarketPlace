package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "ex@test.org",
		Password: "password",
		Active:   true,
	}
}
