package service

import "github.com/VladimirBlinov/MarketPlace/internal/app/store"

type ProductService struct {
	store store.Store
}

func NewProductService(store store.Store) *ProductService{
	return &ProductService{
		store: store,
	}
}

func(ps *ProductService) CreateProduct()