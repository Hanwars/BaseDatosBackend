package models

type Supply struct {
	SupplyID     int32   `json:"supply_id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Stock        int32   `json:"stock"`
	CategoryID   int32   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	Image        string  `json:"image"`
}
