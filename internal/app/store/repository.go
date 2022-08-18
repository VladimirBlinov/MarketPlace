package store

import "github.com/VladimirBlinov/MarketPlace/internal/app/model"

type UserRepo interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type ProductRepo interface {
	Create(*model.Product) error
}
