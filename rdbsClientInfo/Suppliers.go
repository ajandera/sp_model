package rdbsClientInfo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Suppliers struct {
	gorm.Model
	Id         uuid.UUID `gorm:"primary_key; unique"`
	StoreRefer string
	Store      Stores `gorm:"foreignKey:StoreRefer"`
	Name       string
	Street     string
	City       string
	Zip        string
	Country    string
	Email      string
	Phone      string
	Person     string
	Template   string
	Subject    string
}

func (storeWages *Suppliers) BeforeCreate(db *gorm.DB) error {
	storeWages.Id = uuid.New()
	return nil
}
