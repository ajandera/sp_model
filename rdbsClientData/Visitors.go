package rdbsClientData

import (
	"model/rdbsClientInfo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visitors struct {
	gorm.Model
	Id          string `gorm:"primary_key; unique"`
	Ip          string
	StoreId     string
	Url         string
	ProductCode string
	Header      string
	Tag         string
	Store       rdbsClientInfo.Stores `gorm:"foreignKey:Order"`
}

func (visitor *Visitors) BeforeCreate(db *gorm.DB) error {
	visitor.Id = uuid.New().String()
	return nil
}
