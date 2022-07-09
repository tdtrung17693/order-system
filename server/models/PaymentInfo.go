package models

type PaymentMethod struct {
	ID   string `json:"id" gorm:"primarykey,type:varchar(20)"`
	Name string `json:"name"`
}
