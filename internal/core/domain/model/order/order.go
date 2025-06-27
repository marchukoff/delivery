package order

import (
	"errors"
	"slices"

	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var ErrOrderNotInitialized = errors.New("order not initialized")

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		id:       orderID,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func RestoreOrder(id uuid.UUID, courier *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		id:        id,
		courierID: courier,
		location:  location,
		volume:    volume,
		status:    status,
	}
}

func (o *Order) Equals(other *Order) bool {
	ids := []uuid.UUID{o.ID(), other.ID()}
	if slices.Contains(ids, uuid.Nil) {
		return false
	}
	return slices.Contains(ids[1:], ids[0])
}

func (o *Order) ID() uuid.UUID {
	if o == nil {
		return uuid.Nil
	}
	return o.id
}

func (o *Order) CourierID() *uuid.UUID {
	if o == nil {
		return nil
	}
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	if o == nil {
		return kernel.Location{}
	}
	return o.location
}

func (o *Order) Volume() int {
	if o == nil {
		return 0
	}
	return o.volume
}

func (o *Order) Status() Status {
	if o == nil {
		return StatusEmpty
	}
	return o.status
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if o == nil {
		return ErrOrderNotInitialized
	}
	if courierID == uuid.Nil {
		return errs.NewValueIsRequiredError("courierID")
	}
	// ? no reassign
	if !o.status.Equals(StatusCreated) {
		return errs.NewExpectationFailedError("status", o.Status(), StatusCreated)
	}

	o.courierID = &courierID
	o.status = StatusAssigned

	return nil
}

func (o *Order) Complete() error {
	if o == nil {
		return ErrOrderNotInitialized
	}
	if !o.status.Equals(StatusAssigned) {
		return errs.NewExpectationFailedError("status", o.Status(), StatusAssigned)
	}

	o.status = StatusCompleted
	return nil
}
