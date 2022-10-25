package store

import "github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/model"

type UserRepo interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type ProductRepo interface {
	Create(*model.Product, *model.MarketPlaceItemsList) error
	Update(*model.Product, *model.MarketPlaceItemsList) error
	FindByUserId(int) ([]*model.Product, error)
	GetProductById(int) (*model.Product, error)
	GetCategories() ([]*model.Category, error)
	CreateCategory(*model.Category) error
	CreateMaterial(*model.Material) error
	GetMaterials() ([]*model.Material, error)
	Delete(int, int) error
}
