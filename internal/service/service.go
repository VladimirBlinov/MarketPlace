package service

import "github.com/VladimirBlinov/MarketPlace/internal/app/store"

type Service struct {
	ProductService *ProductService
}

func NewService(store store.Store) *Service {
	ProductService := NewProductService(store)
	return &Service{
		ProductService: ProductService,
	}
}
