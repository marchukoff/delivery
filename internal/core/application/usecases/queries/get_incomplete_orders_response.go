package queries

import (
	"github.com/google/uuid"
)

type GetIncompleteOrdersResponse struct {
	Orders []Order
}

type Order struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Location Location  `gorm:"embedded;embeddedPrefix:location_"`
}

func (Order) TableName() string { return "orders" }
