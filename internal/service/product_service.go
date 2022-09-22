package service

import (
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
)

type ProductService struct {
	store store.Store
}

type InputProduct struct {
	ProductName    string  `json:"product_name"`
	CategoryID     int     `json:"category_id,string"`
	PiecesInPack   int     `json:"pieces_in_pack,string"`
	MaterialID     int     `json:"material_id,string"`
	Weight         float32 `json:"weight,string"`
	Lenght         float32 `json:"lenght,string"`
	Width          float32 `json:"width,string"`
	Height         float32 `json:"height,string"`
	Description    string  `json:"description"`
	WildberriesSKU int     `json:"wildberries_sku,string"`
	OzonSKU        int     `json:"ozon_sku,string"`
	UserID         int     `json:"user_id,string"`
}

func NewProductService(store store.Store) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (ps *ProductService) CreateProduct(reqProd InputProduct) error {
	p := &model.Product{
		ProductName:    reqProd.ProductName,
		CategoryID:     reqProd.CategoryID,
		PiecesInPack:   reqProd.PiecesInPack,
		MaterialID:     reqProd.MaterialID,
		Weight:         reqProd.Weight,
		Lenght:         reqProd.Lenght,
		Width:          reqProd.Weight,
		Height:         reqProd.Height,
		Description:    reqProd.Description,
		UserID:         reqProd.UserID,
		Active:         true,
		OzonSKU:        reqProd.OzonSKU,
		WildberriesSKU: reqProd.WildberriesSKU,
	}

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
