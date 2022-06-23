package store

import "github.com/VladimirBlinov/MarketPlace/internal/app/model"

type UserRepo struct {
	store *Store
}

func (r *UserRepo) Create(u *model.User) (*model.User, error) {
	err := u.EncryptPasswordBeforeCreate()
	if err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO public.users (email, encryptedpassword, userrole, active) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		2,
		u.Active,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encryptedpassword, userrole, active FROM public.users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.UserRole,
		&u.Active,
	); err != nil {
		return nil, err
	}

	return u, nil
}
