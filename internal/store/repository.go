package store

import "github.com/VladimirBlinov/MarketPlace/internal/model"

type UserRepo interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type ProductRepo interface {
	Create(*model.Product) error
	FindByUserId(int) ([]*model.Product, error)
	GetCategories() ([]*model.Category, error)
	CreateCategory(*model.Category) error
}