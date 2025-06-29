package queries

import (
	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Couriers []Courier
}

type Courier struct {
	ID       uuid.UUID
	Name     string
	Location Location
}
