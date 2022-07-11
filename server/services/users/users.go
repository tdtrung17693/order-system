package users

import (
	"order-system/database"
	"order-system/models"

	"gorm.io/gorm"
)

func FindUserByEmail(email string) (models.User, error) {
	db := database.GetDBInstance()
	o := models.User{}
	res := db.Preload("Cart").Where("email = ?", email).First(&o)

	return o, res.Error
}

func UserExists(email string) (bool, error) {
	db := database.GetDBInstance()
	var exists bool
	err := db.Model(models.User{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).
		Error
	return exists, err
}

func CreateUser(name string, email string, encryptedPass string, role models.UserRole) (models.User, error) {
	db := database.GetDBInstance()

	o := models.User{
		Email:    email,
		Name:     name,
		Password: encryptedPass,
		Role:     role,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&o).Error; err != nil {
			return err
		}

		// every user will have a cart
		// after registering
		c := models.Cart{
			UserID: o.ID,
		}

		if err := tx.Create(&c).Error; err != nil {
			return err
		}

		return nil
	})

	return o, err
}
