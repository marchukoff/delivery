package courier

import (
	"errors"
	"math"
	"slices"
	"strings"

	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var (
	ErrCourierNotInitialized = errors.New("courier not initialized")
	ErrTargetReached         = errors.New("target reached")
	ErrCannotTakeOrder       = errors.New("cannot carrying anymore")
)

type Courier struct {
	id            uuid.UUID
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if speed <= 0 {
		return nil, errs.NewValueIsRequiredError("speed")
	}

	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}

	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: []*StoragePlace{NewBag()},
	}, nil
}

func RestoreCourier(id uuid.UUID, name string, speed int, location kernel.Location, storagePlaces []*StoragePlace) *Courier {
	return &Courier{
		id:            id,
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: storagePlaces,
	}
}

func (c *Courier) Equals(other *Courier) bool {
	ids := []uuid.UUID{c.ID(), other.ID()}
	if slices.Contains(ids, uuid.Nil) {
		return false
	}
	return slices.Contains(ids[1:], ids[0])
}

func (c *Courier) ID() uuid.UUID {
	if c == nil {
		return uuid.Nil
	}
	return c.id
}

func (c *Courier) Name() string {
	if c == nil {
		return ""
	}
	return c.name
}

func (c *Courier) Speed() int {
	if c == nil {
		return 0
	}
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	if c == nil {
		return kernel.Location{}
	}
	return c.location
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	if c == nil {
		return nil
	}
	return c.storagePlaces
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	if c == nil || c.ID() == uuid.Nil {
		return ErrCourierNotInitialized
	}

	sp, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, sp)

	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if c == nil || c.ID() == uuid.Nil {
		return false, ErrCourierNotInitialized
	}

	if order == nil {
		return false, errs.NewValueIsRequiredError("order")
	}

	can := false
	for _, sp := range c.storagePlaces {
		can, err := sp.CanStore(order.Volume())
		if err != nil {
			return can, err
		}
		if can {
			return can, nil
		}
	}

	return can, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if c == nil || c.ID() == uuid.Nil {
		return ErrCourierNotInitialized
	}

	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	can, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if !can {
		return ErrCannotTakeOrder
	}

	for _, sp := range c.storagePlaces {
		can, err = sp.CanStore(order.Volume())
		if err != nil {
			return err
		}

		if can {
			return sp.Store(order.ID(), order.Volume())
		}
	}

	return ErrCannotTakeOrder
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if c == nil || c.ID() == uuid.Nil {
		return ErrCourierNotInitialized
	}

	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	sp, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return err
	}

	err = sp.Clear(order.ID())
	if err != nil {
		return err
	}

	return order.Complete()
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if c == nil || c.ID() == uuid.Nil {
		return 0, ErrCourierNotInitialized
	}

	if target.IsEmpty() {
		return 0, errs.NewValueIsRequiredError("target")
	}

	dst, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}

	return math.Floor(float64(dst) / float64(c.speed)), nil
}

func (c *Courier) Move(target kernel.Location) error {
	if c == nil || c.ID() == uuid.Nil {
		return ErrCourierNotInitialized
	}

	if target.IsEmpty() {
		return errs.NewValueIsRequiredError("target")
	}

	if target.Equals(c.location) {
		return ErrTargetReached
	}

	// TODO: replace with A*
	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}

	c.location = newLocation
	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	if c == nil || c.ID() == uuid.Nil {
		return nil, ErrCourierNotInitialized
	}

	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}

	for _, sp := range c.storagePlaces {
		if sp.OrderID() != nil && *sp.OrderID() == orderID {
			return sp, nil
		}
	}

	return nil, errs.NewObjectNotFoundError("StoragePlace.orderID", orderID)
}
