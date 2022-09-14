package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "ex@test.org",
		Password: "password",
		UserRole: 2,
		Active:   true,
	}
}

func TestAdminUser(t *testing.T) *User {
	return &User{
		Email:    "ex@test.org",
		Password: "password",
		UserRole: 1,
		Active:   true,
	}
}

func TestProduct(t *testing.T) *Product {
	return &Product{
		ProductName:  "Менажница",
		CategoryID:   105,
		PiecesInPack: 1,
		MaterialID:   1,
		Weight:       500,
		Lenght:       200,
		Width:        300,
		Height:       15,
		Description:  "описание",
		UserID:       1,
		Active:       true,
	}
}

func TestCategory(t *testing.T) *Category {
	return &Category{
		CategoryName:     "Менажница Деревянная",
		ParentCategoryID: 104,
		Active:           true,
	}
}
