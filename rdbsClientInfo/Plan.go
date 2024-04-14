package rdbsClientInfo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	Id       uuid.UUID `gorm:"primary_key; unique"`
	Price    float64
	Period   int
	Name     string
	Products int
	Enabled  bool
	Free     bool
	OneTime  bool
}

func (plan *Plan) BeforeCreate(db *gorm.DB) error {
	plan.Id = uuid.New()
	return nil
}
