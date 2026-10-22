package models

type Client struct {
	ClientID     int32  `json:"client_id"`
	Name         string `json:"name"`
	ContactPhone int32  `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
}
