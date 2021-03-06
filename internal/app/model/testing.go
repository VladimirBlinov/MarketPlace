package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "ex@test.org",
		Password: "password",
		Active:   true,
	}
}

func TestProduct(t *testing.T) *Product {
	return &Product{
		ProductName:  "Менажница",
		CategoryID:   1111110000,
		PiecesInPack: 1,
		MaterialID:   1,
		Weight:       500,
		Lenght:       200,
		Width:        300,
		Height:       15,
		Description:  "описание",
		UserID:       1,
		Avtive:       true,
	}
}
