package teststore

import (
	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
)

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

func (r *ProductRepo) FindByUserId(userId int) ([]*model.Product, error) {
	productsList := make([]*model.Product, len(r.products))
	for _, product := range r.products {
		if product.UserID == userId {
			productsList = append(productsList, product)
		}
	}
	if len(productsList) < 1 {
		return nil, store.ErrRecordNotFound
	}

	return productsList, nil
}
