package courierrepo

import (
	"github.com/google/uuid"
)

type CourierDTO struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string
	Speed         int
	Location      LocationDTO        `gorm:"embedded;embeddedPrefix:location_"`
	StoragePlaces []*StoragePlaceDTO `gorm:"foreignKey:CourierID;constraint:OnDelete:CASCADE;"`
}

func (CourierDTO) TableName() string {
	return "couriers"
}

type LocationDTO struct {
	X, Y int
}

type StoragePlaceDTO struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	TotalVolume int
	OrderID     uuid.UUID `gorm:"type:uuid"`
	CourierID   uuid.UUID `gorm:"type:uuid"`
}

func (StoragePlaceDTO) TableName() string {
	return "storage_places"
}
