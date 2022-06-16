package model

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
	Avtive       bool
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
