package models

type SaleDetail struct {
	DetailID   int32   `json:"detail_id"`
	SaleID     int32   `json:"sale_id"`
	SupplyID   int32   `json:"supply_id"`
	SupplyName string  `json:"supply_name"`
	Quantity   int32   `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	Discount   float64 `json:"discount"`
	Subtotal   float64 `json:"subtotal"`
}
