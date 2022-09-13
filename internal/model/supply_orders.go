package model

import (
	"time"
)

type SupplyOrder struct {
	SupplyOrderID          int
	OrderDate              time.Time
	SupplierID             int
	ShippingCostToLogistic float32
	ShippingCostByLogistic float32
	Active                 bool
}

type SupplyOrderProduct struct {
	SupplyOrderProductID int
	SupplyOrderID        int
	ProductID            int
	Quantity             float32
	UnitPrice            float32
	CurrencyID           int
}

type Supplier struct {
	SupplierID            int
	SupplierName          string
	SupplierAddress       string
	SupplierCountryID     int
	SupplierSWIFT         string
	SupplierAccountNumber int
	Active                bool
}

type SupplyOrderStatus struct {
	SupplyOrderStatusID   int
	SupplyOrderStatusName string
}

type Currency struct {
	CurrencyID   int
	CurrencyName string
	CurrencyCode string
}

type PaymentStatus struct {
	PaymentStatusID    int
	PaymentStatus_Name string
}

type Payment struct {
	PaymentID       int
	PaymentDate     time.Time
	PaymentAmount   float32
	CurrencyID      int
	SupplierID      int
	SupplyOrderID   int
	PaymentStatusID int
	Active          bool
}

type SupplyOrderPayment struct {
	SupplyOrderPaymentID int
	SupplyOrderID        int
	PaymentID            int
}

type SupplyOrderAudit struct {
	SupplyOrderAuditID   int
	SupplyOrderID        int
	SupplyOrderAuditDate time.Time
	SupplyOrderStatusID  int
	AuditUserID          int
	Active               bool
}

type PaymentAudit struct {
	PaymentAuditID   int
	PaymentID        int
	PaymentAuditDate time.Time
	PaymentStatusID  int
	AuditUserID      int
	Active           bool
}
