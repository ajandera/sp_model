package rdbsClientData

import (
	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItems struct {
	gorm.Model
	Id          string `gorm:"primary_key; unique"`
	UnitPrice   float64
	Quantity    int8
	ProductCode string
	Order       string
	ProductName string
	OrderEntity rdbsClientInfo.Stores `gorm:"foreignKey:Order"`
}

func (orderItem *OrderItems) BeforeCreate(db *gorm.DB) error {
	orderItem.Id = uuid.New().String()
	return nil
}
