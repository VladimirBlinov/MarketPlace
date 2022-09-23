package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Product struct {
	ProductID      int     `json:"product_id,string"`
	ProductName    string  `json:"product_name"`
	CategoryID     int     `json:"category_id,string"`
	PiecesInPack   int     `json:"pieces_in_pack,string"`
	MaterialID     int     `json:"material_id,string"`
	Weight         float32 `json:"weight,string"`
	Lenght         float32 `json:"lenght,string"`
	Width          float32 `json:"width,string"`
	Height         float32 `json:"height,string"`
	Description    string  `json:"description"`
	UserID         int     `json:"user_id"`
	Active         bool    `json:"active"`
	WildberriesSKU int     `json:"wildberries_sku,string"`
	OzonSKU        int     `json:"ozon_sku,string"`
}

func (p *Product) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.ProductName, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.CategoryID, validation.Required),
		validation.Field(&p.CategoryID, validation.Required, validation.By(checkCategoryID(1))),
		validation.Field(&p.MaterialID, validation.Required),
		validation.Field(&p.UserID, validation.Required),
		// validation.Field(&p.Active, validation.Required),
	)
}

type MarketPlaceItem struct {
	MarketPlaceItemID int
	ProductID         int
	ItemName          string
	MarketPlaceID     int
	SKU               int
	UserID            int
	Active            bool
}

func (mpi *MarketPlaceItem) ValidateMarketPlaceItem() error {
	return validation.ValidateStruct(
		mpi,
		validation.Field(&mpi.MarketPlaceID, validation.Required),
		validation.Field(&mpi.UserID, validation.Required),
	)
}

type MarketPlaceItemsList struct {
	MPIList []*MarketPlaceItem
}

func (mpil *MarketPlaceItemsList) GetMPIList(p *Product) {
	if p.OzonSKU != 0 {
		mpiOzon := &MarketPlaceItem{
			ItemName: p.ProductName,
			UserID:   p.UserID,
			Active:   true,
		}
		mpiOzon.MarketPlaceID = 1
		mpiOzon.SKU = p.OzonSKU
		mpil.MPIList = append(mpil.MPIList, mpiOzon)
	}

	if p.WildberriesSKU != 0 {
		mpiWb := &MarketPlaceItem{
			ItemName: p.ProductName,
			UserID:   p.UserID,
			Active:   true,
		}
		mpiWb.MarketPlaceID = 2
		mpiWb.SKU = p.WildberriesSKU
		mpil.MPIList = append(mpil.MPIList, mpiWb)
	}
}

func (mpi *MarketPlaceItemsList) ValidateMarketPlaceItems() error {
	for _, mpi := range mpi.MPIList {
		err := mpi.ValidateMarketPlaceItem()
		if err != nil {
			return err
		}
	}

	return nil
}

type Category struct {
	CategoryID       int    `json:"category_id"`
	CategoryName     string `json:"category_name"`
	ParentCategoryID int    `json:"parent_category_id"`
	Active           bool   `json:"active"`
}

func (c *Category) ValidateCategory() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.CategoryName, validation.Required, validation.Length(3, 200)),
	)
}

type MarketPlace struct {
	MarketPlaceID   int
	MarketPlaceName string
	Active          bool
}

type Material struct {
	MaterialID   int    `json:"material_id"`
	MaterialName string `json:"material_name"`
	Active       bool   `json:"active"`
}

func (m *Material) ValidateMaterial() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.MaterialName, validation.Required, validation.Length(3, 200)))
}
