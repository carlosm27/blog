package model

import (
	"gorm.io/gorm"
)

type Expenses struct {
	gorm.Model
	Amount      float64 `json:"amount" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Category    string  `json:"category" gorm:"not null"`
}
