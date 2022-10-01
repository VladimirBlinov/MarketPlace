package service

import (
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductService struct {
	store store.Store
}

func NewProductService(store store.Store) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (ps *ProductService) CreateProduct(p *model.Product) error {
	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)

	if err := p.Validate(); err != nil {
		return err
	}
	if err := mpiList.ValidateMarketPlaceItems(); err != nil {
		return err
	}
	if err := ps.store.Product().Create(p, mpiList); err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) GetProductById(id int) (*model.Product, error) {
	product, err := ps.store.Product().GetProductById(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductService) GetProductsByUserId(userId int) ([]*model.Product, error) {
	products, err := ps.store.Product().FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) GetProductCategories() ([]*model.Category, error) {
	categories, err := ps.store.Product().GetCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (ps *ProductService) GetProductMaterials() ([]*model.Material, error) {
	materials, err := ps.store.Product().GetMaterials()
	if err != nil {
		return nil, err
	}

	return materials, nil
}
