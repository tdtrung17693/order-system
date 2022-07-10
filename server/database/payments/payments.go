package payments

import (
	"order-system/database"
	"order-system/handlers/dto"
	"order-system/models"
)

func FindAllPaymentMethod() ([]dto.PaymentMethodDto, error) {
	db := database.GetDBInstance()
	p := []models.PaymentMethod{}

	err := db.Find(&p).Error
	if err != nil {
		return nil, err
	}

	result := []dto.PaymentMethodDto{}
	for _, pm := range p {
		result = append(result, dto.PaymentMethodDto{
			Id:   pm.ID,
			Name: pm.Name,
		})
	}

	return result, nil
}
