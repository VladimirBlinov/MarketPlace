package model

import (
	"time"
)

type SupplyOrder struct {
	SupplyOrderID          int
	OrderDate              time.Time
	ProductID              int
	SupplierID             int
	Quantity               float32
	UnitPrice              float32
	CurrencyID             int
	ShippingCostToLogistic float32
	ShippingCostByLogistic float32
	ShippingDate           time.Time
	ArrivingDate           time.Time
	PaymentID              int
	IsPaid                 bool
	Avtive                 bool
}

type Supplier struct {
	SupplierID            int
	SupplierName          int
	SupplierAddress       string
	SupplierCountryID     int
	SupplierSWIFT         string
	SupplierAccountNumber int
	Avtive                bool
}
