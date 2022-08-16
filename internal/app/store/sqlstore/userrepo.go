package sqlstore

import (
	"database/sql"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
)

type UserRepo struct {
	store *Store
}

func (r *UserRepo) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	err := u.EncryptPasswordBeforeCreate()
	if err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO public.users (email, encryptedpassword, userrole, active) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.UserRole,
		u.Active,
	).Scan(&u.ID)
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
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}