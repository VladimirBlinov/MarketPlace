package teststore

import (
	"github.com/VladimirBlinov/MarketPlace/Marketplace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/Marketplace/internal/store"
	_ "github.com/lib/pq"
)

// Store
type Store struct {
	userRepo    *UserRepo
	ProductRepo *ProductRepo
}

// Store constructor
func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepo{
		store: s,
		users: make(map[int]*model.User),
	}
	return s.userRepo
}

func (s *Store) Product() store.ProductRepo {
	if s.ProductRepo != nil {
		return s.ProductRepo
	}

	s.ProductRepo = &ProductRepo{
		store:            s,
		Products:         make(map[int]*model.Product),
		categories:       make(map[int]*model.Category),
		materials:        make(map[int]*model.Material),
		marketPlaceItems: make(map[int]*model.MarketPlaceItem),
	}
	return s.ProductRepo
}
