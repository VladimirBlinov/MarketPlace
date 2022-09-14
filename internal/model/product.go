package model

import validation "github.com/go-ozzo/ozzo-validation"

type Product struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CategoryID   int     `json:"category_id"`
	PiecesInPack int     `json:"pieces_in_pack"`
	MaterialID   int     `json:"material_id"`
	Weight       float32 `json:"weight"`
	Lenght       float32 `json:"lenght"`
	Width        float32 `json:"width"`
	Height       float32 `json:"height"`
	Description  string  `json:"description"`
	UserID       int     `json:"user_id"`
	Active       bool    `json:"active"`
}

func (p *Product) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.ProductName, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.CategoryID, validation.Required),
		//validation.Field(&p.CategoryID, validation.Required, validation.By(checkCategoryID(1))),
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
	Barcode           int
	Active            bool
}

type Category struct {
	CategoryID       int
	CategoryName     string
	ParentCategoryID int
	Active           bool
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