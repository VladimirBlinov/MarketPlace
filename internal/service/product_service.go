package service

import (
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductService struct {
	store store.Store
}

type RequestProduct struct {
	ProductName    string  `json:"product_name"`
	CategoryID     int     `json:"category_id,string"`
	PiecesInPack   int     `json:"pieces_in_pack,string"`
	MaterialID     int     `json:"material_id,string"`
	Weight         float32 `json:"weight,string"`
	Lenght         float32 `json:"lenght,string"`
	Width          float32 `json:"width,string"`
	Height         float32 `json:"height,string"`
	Description    string  `json:"description"`
	MarketPlaceID  int     `json:"marketplace_id"`
	MarketPlaceSKU int     `json:"marketplace_sku"`
}

func NewProductService(store store.Store) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (ps *ProductService) CreateProduct(reqProd RequestProduct, userId int) error {
	p := &model.Product{
		ProductName:  reqProd.ProductName,
		CategoryID:   reqProd.CategoryID,
		PiecesInPack: reqProd.PiecesInPack,
		MaterialID:   reqProd.MaterialID,
		Weight:       reqProd.Weight,
		Lenght:       reqProd.Lenght,
		Width:        reqProd.Weight,
		Height:       reqProd.Height,
		Description:  reqProd.Description,
		UserID:       userId,
		Active:       true,
	}

	if err := ps.store.Product().Create(p); err != nil {
		return err
	}

	return nil
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
