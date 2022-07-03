package model

import (
	"gorm.io/gorm"
)

type Expenses struct {
	gorm.Model
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}
