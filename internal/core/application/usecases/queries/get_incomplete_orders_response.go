package queries

import (
	"github.com/google/uuid"
)

type GetIncompleteOrdersResponse struct {
	Orders []Order
}

type Order struct {
	ID       uuid.UUID
	Location Location
}
