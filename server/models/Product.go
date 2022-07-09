package models

type Product struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
	Vendor      User   `json:"vendor"`
	VendorID    uint   `json:"vendorId"`
	Unit        string `json:"unit"`
}
