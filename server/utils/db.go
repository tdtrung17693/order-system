package utils

import (
	"sync"

	"gorm.io/gorm/schema"
)

func GetModelTableName(model interface{}) string {
	s, _ := schema.Parse(&model, &sync.Map{}, schema.NamingStrategy{})
	return s.Table
}
