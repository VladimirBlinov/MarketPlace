package store

import "github.com/VladimirBlinov/MarketPlace/internal/app/model"

type ProductRepo struct {
	store *Store
}

func (r *ProductRepo) Create(p *model.Product) (*model.Product, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO public.product (product_name, category_id, pieces_in_pack, material_id, weight_gr, lenght_mm, width_mm, height_mm, product_description, user_id, avtive) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING product_id",
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
		p.Avtive,
	).Scan(&p.ProductID); err != nil {
		return nil, err
	}

	return p, nil
}
