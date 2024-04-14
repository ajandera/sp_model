package rdbsClientInfo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreWeights struct {
	gorm.Model
	Id                 uuid.UUID `gorm:"primary_key; unique"`
	StoreRefer         string
	Store              Stores `gorm:"foreignKey:StoreRefer"`
	Name               string
	Beta               float64
	Gama               float64
	Delta              float64
	A                  float64
	B                  float64
	C                  float64
	D                  float64
	E                  float64
	ProbabilityWeights string
	Shift              int
	LongShift          int
}

func (storeWages *StoreWeights) BeforeCreate(db *gorm.DB) error {
	storeWages.Id = uuid.New()
	return nil
}
