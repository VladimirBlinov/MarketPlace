package teststore

import "github.com/VladimirBlinov/MarketPlace/internal/app/model"

type ProductRepo struct {
	store    *Store
	products map[int]*model.Product
}

func (r *ProductRepo) Create(p *model.Product) error {
	if err := p.Validate(); err != nil {
		return nil
	}

	p.ProductID = len(r.products) + 1
	r.products[p.ProductID] = p

	return nil
}
