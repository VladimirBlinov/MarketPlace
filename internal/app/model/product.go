package model

import validation "github.com/go-ozzo/ozzo-validation"

type Product struct {
	ProductID    int
	ProductName  string
	CategoryID   int
	PiecesInPack int
	MaterialID   int
	Weight       float32
	Lenght       float32
	Width        float32
	Height       float32
	Description  string
	UserID       int
	Avtive       bool
}

func (p *Product) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.ProductName, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.CategoryID, validation.Required),
		validation.Field(&p.PiecesInPack, validation.Required),
		validation.Field(&p.MaterialID, validation.Required),
		validation.Field(&p.UserID, validation.Required),
		validation.Field(&p.Avtive, validation.Required),
	)
}

type MarketPlaceItem struct {
	MarketPlaceItemID int
	ProductID         int
	ItemName          string
	MarketPlaceID     int
	Barcode           int
	Avtive            bool
}

type Category struct {
	CategoryID       int
	CategoryName     string
	ParentCategoryID int
	Active           bool
}

type MarketPlace struct {
	MarketPlaceID   int
	MarketPlaceName string
	Active          bool
}
