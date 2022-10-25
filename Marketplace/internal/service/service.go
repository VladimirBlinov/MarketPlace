package service

import "github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store"

type Service struct {
	ProductService *ProductService
	AuthService    *AuthService
}

func NewService(store store.Store) *Service {
	ProductService := NewProductService(store)
	AuthService := NewAuthService(store)
	return &Service{
		ProductService: ProductService,
		AuthService:    AuthService,
	}
}
