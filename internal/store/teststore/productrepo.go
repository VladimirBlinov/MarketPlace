package teststore

import (
	"errors"
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductRepo struct {
	store            *Store
	products         map[int]*model.Product
	categories       map[int]*model.Category
	materials        map[int]*model.Material
	marketPlaceItems map[int]*model.MarketPlaceItem
}

func (r *ProductRepo) Create(p *model.Product, mpiList *model.MarketPlaceItemsList) error {
	p.ProductID = len(r.products) + 1
	r.products[p.ProductID] = p

	for _, mpi := range mpiList.MPIList {
		mpi.ProductID = p.ProductID
		mpi.MarketPlaceItemID = len(r.marketPlaceItems) + 1
		r.marketPlaceItems[mpi.MarketPlaceItemID] = mpi
	}

	return nil
}

func (r *ProductRepo) GetProductById(productId int) (*model.Product, error) {
	for _, product := range r.products {
		if product.ProductID == productId {
			for _, mpi := range r.marketPlaceItems {
				GetProductIdMarketplaceItem(product, mpi)
			}
			return product, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func GetProductIdMarketplaceItem(p *model.Product, mpi *model.MarketPlaceItem) {
	if mpi.ProductID == p.ProductID {
		switch mpi.MarketPlaceID {
		case 1:
			p.OzonSKU = mpi.SKU
		case 2:
			p.WildberriesSKU = mpi.SKU
		default:
			break
		}
	}
}

func (r *ProductRepo) FindByUserId(userId int) ([]*model.Product, error) {
	productsList := make([]*model.Product, 0)
	for _, product := range r.products {
		if product.UserID == userId {
			for _, mpi := range r.marketPlaceItems {
				if mpi.ProductID == product.ProductID {
					switch mpi.MarketPlaceID {
					case 1:
						product.OzonSKU = mpi.SKU
					case 2:
						product.WildberriesSKU = mpi.SKU
					default:
						return nil, errors.New("wrong MarketPlaceID")
					}
				}
			}
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
		return err
	}

	c.CategoryID = len(r.categories) + 1
	r.categories[c.CategoryID] = c

	return nil
}

func (r *ProductRepo) CreateMaterial(m *model.Material) error {
	if err := m.ValidateMaterial(); err != nil {
		return err
	}

	m.MaterialID = len(r.materials) + 1
	r.materials[m.MaterialID] = m

	return nil
}

func (r *ProductRepo) GetMaterials() ([]*model.Material, error) {
	materials := make([]*model.Material, 0)

	for _, material := range r.materials {
		materials = append(materials, material)
	}

	if len(materials) < 1 {
		return nil, store.ErrRecordNotFound
	}

	return materials, nil
}
