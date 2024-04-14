package rdbsClientData

import (
	"model/rdbsClientInfo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	Id              string `gorm:"primary_key; unique"`
	Amount          float64
	Currency        string
	StoreId         string
	ExternalOrderId string
	Tag             string
	Store           rdbsClientInfo.Stores `gorm:"foreignKey:Order"`
}

func (order *Orders) BeforeCreate(db *gorm.DB) error {
	order.Id = uuid.New().String()
	return nil
}
