package teststore

import (
	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
	_ "github.com/lib/pq"
)

// Store
type Store struct {
	userRepo    *UserRepo
	productRepo *ProductRepo
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
	if s.productRepo != nil {
		return s.productRepo
	}

	s.productRepo = &ProductRepo{
		store: s,
	}
	return s.productRepo
}
