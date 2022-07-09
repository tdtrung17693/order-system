package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseWithPrimaryKey struct {
	ID uint `json:"id" gorm:"primarykey"`
}

type BaseWithAudit struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BaseWithSoftDelete struct {
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Base struct {
	BaseWithPrimaryKey
	BaseWithAudit
	BaseWithSoftDelete
}
