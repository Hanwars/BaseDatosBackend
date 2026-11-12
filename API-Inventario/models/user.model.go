package models

type User struct {
	UserID   int32
	RoleID   int32
	Name     string
	Password string
	Email    string
}

type UserAdmin struct {
	UserID   int32   `json:"user_id"`
	RoleID   int32   `json:"role_id"`
	RoleName string  `json:"role_name"`
	Name     string  `json:"name"`
	Email    *string `json:"email,omitempty"`
}
