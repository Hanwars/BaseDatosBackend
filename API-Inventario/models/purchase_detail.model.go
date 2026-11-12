package models

type PurchaseDetail struct {
	DetailID   int32   `json:"detail_id"`
	PurchaseID int32   `json:"purchase_id"`
	SupplyID   int32   `json:"supply_id"`
	SupplyName string  `json:"supply_name"`
	Quantity   int32   `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	Discount   float64 `json:"discount"`
	Subtotal   float64 `json:"subtotal"`
}
