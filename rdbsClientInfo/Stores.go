package rdbsClientInfo

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stores struct {
	gorm.Model
	Id                         uuid.UUID `gorm:"primary_key; unique"`
	CountryCode                string
	LastPrediction             time.Time
	Url                        string
	MaximalProductPrice        float64
	MinimalProductPrice        float64
	ActualStorePower           float64
	ActualCustomerSatisfaction float64
	PerceivedValue             float64
	Code                       string
	AccountRefer               string
	Account                    Accounts `gorm:"foreignKey:AccountRefer"`
	ProductSell                int
	Offline                    bool
	ShoptetId                  string
	ShoptetAccessToken         string
	XmlFeed                    string
	Window                     int8
}

func (stores *Stores) BeforeCreate(db *gorm.DB) error {
	stores.Id = uuid.New()
	return nil
}
