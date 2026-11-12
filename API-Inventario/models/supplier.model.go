package models

type Supplier struct {
	SupplierID   int32   `json:"supplier_id"`
	CompanyName  string  `json:"company_name"`
	CompanyPhone int32   `json:"company_phone"`
	CompanyEmail string  `json:"company_email"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *int32  `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
}
