package model

type Product struct {
	ProductID    int
	Name         string
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
	Name              string
	MarketPlaceID     int
	Barcode           int
	Avtive            bool
}
