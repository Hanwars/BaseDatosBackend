package models

import "time"

type Sale struct {
	SaleID     int32     `json:"sale_id"`
	UserID     int32     `json:"user_id"`
	UserName   string    `json:"user_name"`
	ClientID   int32     `json:"client_id"`
	ClientName string    `json:"client_name"`
	Date       time.Time `json:"date"`
	Subtotal   float64   `json:"subtotal"`
	Discount   float64   `json:"discount"`
	IVA        float64   `json:"iva"`
	Total      float64   `json:"total"`
}
