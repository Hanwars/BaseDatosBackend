package models

import "time"

type SalesSummary struct {
	ClientID     int32      `json:"client_id"`
	ClientName   string     `json:"client_name"`
	SalesCount   int64      `json:"sales_count"`
	TotalSales   float64    `json:"total_sales"`
	LastSaleDate *time.Time `json:"last_sale_date,omitempty"`
}

type InventoryReportItem struct {
	SupplyID     int32   `json:"supply_id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Stock        int32   `json:"stock"`
	CategoryID   int32   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	Image        string  `json:"image"`
	StockStatus  string  `json:"stock_status"`
}

type DiscountMovement struct {
	MovementID int32   `json:"movement_id"`
	Subtotal   float64 `json:"subtotal"`
	Discount   float64 `json:"discount"`
	Type       string  `json:"type"`
}

type AuditRecord struct {
	AuditTable      string    `json:"audit_table"`
	AuditID         int64     `json:"audit_id"`
	AffectedTable   string    `json:"affected_table"`
	Operation       string    `json:"operation"`
	RecordID        int64     `json:"record_id"`
	DateTime        time.Time `json:"date_time"`
	SQLLogin        *string   `json:"sql_login,omitempty"`
	ApplicationUser *int32    `json:"application_user_id,omitempty"`
	PreviousValues  *string   `json:"previous_values,omitempty"`
	NewValues       *string   `json:"new_values,omitempty"`
	PasswordChanged *bool     `json:"password_changed,omitempty"`
}

type BackupRecord struct {
	AuditID         int64     `json:"audit_id"`
	FileName        string    `json:"file_name"`
	Path            string    `json:"path"`
	CreatedAt       time.Time `json:"created_at"`
	Status          string    `json:"status"`
	Details         *string   `json:"details,omitempty"`
	ApplicationUser *int32    `json:"application_user_id,omitempty"`
	SQLLogin        *string   `json:"sql_login,omitempty"`
}
