package commands

import (
	"strings"

	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type AddStoragePlaceCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int
	valid       bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlaceCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errs.NewValueIsRequiredError("courierID")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return AddStoragePlaceCommand{}, errs.NewValueIsRequiredError("name")
	}

	if totalVolume <= 0 {
		return AddStoragePlaceCommand{}, errs.NewValueIsRequiredError("totalVolume")
	}

	return AddStoragePlaceCommand{
		courierID:   courierID,
		name:        name,
		totalVolume: totalVolume,
		valid:       true,
	}, nil
}

func (c AddStoragePlaceCommand) CourierID() uuid.UUID { return c.courierID }

func (c AddStoragePlaceCommand) Name() string { return c.name }

func (c AddStoragePlaceCommand) TotalVolume() int { return c.totalVolume }

func (c AddStoragePlaceCommand) IsValid() bool { return c.valid }
