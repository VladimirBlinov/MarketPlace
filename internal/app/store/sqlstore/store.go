package sqlstore

import (
	"database/sql"

	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
	_ "github.com/lib/pq"
)

// Store
type Store struct {
	db          *sql.DB
	userRepo    *UserRepo
	productRepo *ProductRepo
}

// Store constructor
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepo{
		store: s,
	}
	return s.userRepo
}

func (s *Store) Product() *ProductRepo {
	if s.productRepo != nil {
		return s.productRepo
	}

	s.productRepo = &ProductRepo{
		store: s,
	}
	return s.productRepo
}
