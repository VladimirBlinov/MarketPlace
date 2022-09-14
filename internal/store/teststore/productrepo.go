package teststore

import (
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductRepo struct {
	store      *Store
	products   map[int]*model.Product
	categories map[int]*model.Category
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
	productsList := make([]*model.Product, 0)
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

func (r *ProductRepo) GetCategories() ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
	for _, category := range r.categories {
		categories = append(categories, category)
	}

	if len(r.categories) < 1 {
		return nil, store.ErrRecordNotFound
	}

	return categories, nil
}

func (r *ProductRepo) CreateCategory(c *model.Category) error {
	if err := c.ValidateCategory(); err != nil {
		return nil
	}

	c.CategoryID = len(r.categories) + 1
	r.categories[c.CategoryID] = c

	return nil
}
