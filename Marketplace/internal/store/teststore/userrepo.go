package teststore

import (
	"github.com/VladimirBlinov/MarketPlace/Marketplace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/Marketplace/internal/store"
)

type UserRepo struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepo) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	err := u.EncryptPasswordBeforeCreate()
	if err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound

}

func (r *UserRepo) FindById(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
