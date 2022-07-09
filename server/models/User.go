package models

type UserRole int32

const (
	RegularUser UserRole = iota
	Vendor
)

type User struct {
	Base
	Name     string   `json:"name"`
	Role     UserRole `json:"role"`
	Email    string   `json:"email"`
	Password string   `json:"-"`
	Cart     Cart     `json:"cart"`
}
