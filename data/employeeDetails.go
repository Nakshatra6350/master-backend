package data

import (
	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	FirstName string  `gorm:"type:varchar(100)"`
	LastName  string  `gorm:"type:varchar(100)"`
	Email     string  `gorm:"type:varchar(100);unique_index"`
	Salary    float64 `gorm:"type:decimal(10,2)"`
}
