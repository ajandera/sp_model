package rdbsClientData

import (
	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	Id          string `gorm:"primary_key; unique"`
	Quantity    int8
	ProductCode string
	Name        string
	StoreId     string
	Store       rdbsClientInfo.Stores `gorm:"foreignKey:Order"`
}

func (products *Products) BeforeCreate(db *gorm.DB) error {
	products.Id = uuid.New().String()
	return nil
}
