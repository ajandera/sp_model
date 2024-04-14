package rdbsClientInfo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	Id            uuid.UUID `gorm:"primary_key; unique"`
	AccountRefer  string
	Account       Accounts `gorm:"foreignKey:AccountRefer"`
	StoreRefer    string
	Store         Stores `gorm:"foreignKey:StoreRefer"`
	PlanRefer     string
	Plan          Plan `gorm:"foreignKey:PlanRefer"`
	Amount        float64
	Paid          bool
	Name          string
	Email         string
	Street        string
	City          string
	Zip           string
	CountryCode   string
	CompanyNumber string
	VatNumber     string
	Number        string
}

func (orders *Orders) BeforeCreate(db *gorm.DB) error {
	orders.Id = uuid.New()
	return nil
}
