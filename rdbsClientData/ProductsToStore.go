package rdbsClientData

import (
	"time"

	"model/rdbsClientInfo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductsToStore struct {
	gorm.Model
	Id          string `gorm:"primary_key; unique"`
	Quantity    int8
	ProductCode string
	DateToNeed  time.Time
	DateToOrder time.Time
	StoreId     string
	Store       rdbsClientInfo.Stores `gorm:"foreignKey:Order"`
}

func (product *ProductsToStore) BeforeCreate(db *gorm.DB) error {
	product.Id = uuid.New().String()
	return nil
}
