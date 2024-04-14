package rdbsClientInfo

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Accounts struct {
	gorm.Model
	Id                     uuid.UUID `gorm:"primary_key; unique"`
	Name                   string
	Email                  string
	Street                 string
	City                   string
	Zip                    string
	CountryCode            string
	CompanyNumber          string
	VatNumber              string
	Password               string
	RestoreToken           string
	Parent                 string
	Role                   string
	ValidTokenTo           time.Time
	NewsletterConfirmation time.Time
	Newsletter             bool
}

func (account *Accounts) BeforeCreate(db *gorm.DB) error {
	account.Id = uuid.New()
	return nil
}
