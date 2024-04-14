package rdbsClientData

import (
	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VisitorsOffline struct {
	gorm.Model
	Id      string `gorm:"primary_key; unique"`
	Info    string
	StoreId string
	Store   rdbsClientInfo.Stores `gorm:"foreignKey:Order"`
}

func (visitorOffline *VisitorsOffline) BeforeCreate(db *gorm.DB) error {
	visitorOffline.Id = uuid.New().String()
	return nil
}
