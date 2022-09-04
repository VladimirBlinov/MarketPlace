package sqlstore

import (
	"database/sql"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
)

type ProductRepo struct {
	store *Store
}

func (r *ProductRepo) Create(p *model.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO public.product (product_name, category_id, pieces_in_pack, material_id, weight_gr, lenght_mm, width_mm, height_mm, product_description, user_id, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING product_id",
		p.ProductName,
		p.CategoryID,
		p.PiecesInPack,
		p.MaterialID,
		p.Weight,
		p.Lenght,
		p.Weight,
		p.Height,
		p.Description,
		p.UserID,
		p.Active,
	).Scan(&p.ProductID)
}

func (r *ProductRepo) FindByUserId(userId int) ([]*model.Product, error) {

	var products []*model.Product

	rows, err := r.store.db.Query(
		"SELECT product_id, product_name, category_id, pieces_in_pack ,material_id, weight_gr, lenght_mm, width_mm, height_mm, product_description, user_id, active, FROM public.Product WHERE active = 1 and userid = $1",
		userId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		p := &model.Product{}
		if err := rows.Scan(
			&p.ProductID,
			&p.ProductName,
			&p.CategoryID,
			&p.PiecesInPack,
			&p.MaterialID,
			&p.Weight,
			&p.Lenght,
			&p.Width,
			&p.Height,
			&p.Description,
			&p.UserID,
			&p.Active,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
