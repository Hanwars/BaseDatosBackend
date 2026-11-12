package models

import "time"

type Purchase struct {
	PurchaseID   int32     `json:"purchase_id"`
	UserID       int32     `json:"user_id"`
	UserName     string    `json:"user_name"`
	SupplierID   int32     `json:"supplier_id"`
	SupplierName string    `json:"supplier_name"`
	Date         time.Time `json:"date"`
	Subtotal     float64   `json:"subtotal"`
	Discount     float64   `json:"discount"`
	IVA          float64   `json:"iva"`
	Total        float64   `json:"total"`
}
