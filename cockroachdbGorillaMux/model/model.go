package model

import (
	"gorm.io/gorm"
)

type Expenses struct {
	gorm.Model
	ID          uint    `json:"id"`
	IdExpense   string  `json:"id_expense" gorm:"type:uuid;default:uuid_generate_v4()"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}
