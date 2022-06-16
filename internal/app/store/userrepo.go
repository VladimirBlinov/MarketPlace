package store

import "github.com/VladimirBlinov/MarketPlace/internal/app/model"

type UserRepo struct {
	store *Store
}

func (r *UserRepo) Create(u *model.User) (*model.User, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO public.users (email, encryptedpassword, userrole) VALUES ($1, $2, $3) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		2,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	return nil, nil
}
