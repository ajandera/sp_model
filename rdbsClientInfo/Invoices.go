package rdbsClientInfo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Invoices struct {
	gorm.Model
	Id         uuid.UUID `gorm:"primary_key; unique"`
	StoreRefer string
	Store      Stores `gorm:"foreignKey:StoreRefer"`
	DueDate    time.Time
	Amount     float64
	Currency   string
}

func (invoice *Invoices) BeforeCreate(db *gorm.DB) error {
	invoice.Id = uuid.New()
	return nil
}
