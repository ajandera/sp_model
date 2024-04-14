package rdbsClientInfo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OpenData struct {
	gorm.Model
	Id                   uuid.UUID `gorm:"primary_key; unique"`
	StorePower           float64
	CustomerSatisfaction float64
	MaximalProductPrice  float64
	MinimalProductPrice  float64
	PerceivedValue       float64
	StoreRefer           string
	Store                Stores `gorm:"foreignKey:StoreRefer"`
}

func (openData *OpenData) BeforeCreate(db *gorm.DB) error {
	openData.Id = uuid.New()
	return nil
}
