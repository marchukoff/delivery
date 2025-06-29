package commands

import (
	"strings"

	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string
	volume  int
	valid   bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int) (CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("orderID")
	}

	street = strings.TrimSpace(street)
	if street == "" {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("street")
	}

	if volume <= 0 {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("volume")
	}

	return CreateOrderCommand{
		orderID: orderID,
		street:  street,
		volume:  volume,
		valid:   true,
	}, nil
}

func (c CreateOrderCommand) OrderID() uuid.UUID { return c.orderID }

func (c CreateOrderCommand) Street() string { return c.street }

func (c CreateOrderCommand) Volume() int { return c.volume }

func (c CreateOrderCommand) IsValid() bool { return c.valid }
