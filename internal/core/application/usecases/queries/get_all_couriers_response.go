package queries

import (
	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Couriers []Courier
}

type Courier struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	Location Location `gorm:"embedded;embeddedPrefix:location_"`
}

func (Courier) TableName() string { return "couriers" }
