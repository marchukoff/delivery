package courier

import (
	"errors"
	"strings"

	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

var (
	ErrStoragePlaceNotInitialized = errors.New("storage place not initialized")
	ErrStoragePlaceIsOccupied     = errors.New("storage place is occupied")
	ErrTotalVolumeMustPositive    = errors.New("totalVolume should be greater than 0")
	ErrVolumeMustPositive         = errors.New("volume should be greater than 0")
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if totalVolume <= 0 {
		return nil, errs.NewValueIsInvalidErrorWithCause("totalVolume", ErrTotalVolumeMustPositive)
	}

	return &StoragePlace{id: uuid.New(), name: name, totalVolume: totalVolume}, nil
}

func NewBag() *StoragePlace {
	const (
		name   = "Сумка"
		volume = 10
	)

	bag, err := NewStoragePlace(name, volume)
	if err != nil {
		panic(err)
	}

	return bag
}

func (sp *StoragePlace) Equals(other *StoragePlace) bool {
	if sp == nil || other == nil || sp.id == uuid.Nil || other.id == uuid.Nil {
		return false
	}

	return sp.id == other.id
}

func (sp *StoragePlace) ID() uuid.UUID {
	if sp == nil {
		return uuid.Nil
	}

	return sp.id
}

func (sp *StoragePlace) Name() string {
	if sp == nil {
		return ""
	}

	return sp.name
}

func (sp *StoragePlace) TotalVolume() int {
	if sp == nil {
		return 0
	}

	return sp.totalVolume
}

func (sp *StoragePlace) OrderID() *uuid.UUID {
	if sp == nil {
		return nil
	}

	return sp.orderID
}

func (sp *StoragePlace) CanStore(volume int) (bool, error) {
	if sp == nil {
		return false, ErrStoragePlaceNotInitialized
	}

	if volume <= 0 {
		return false, errs.NewValueIsInvalidErrorWithCause("volume", ErrVolumeMustPositive)
	}

	return !sp.IsOccupied() && sp.totalVolume >= volume, nil
}

func (sp *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if sp == nil {
		return ErrStoragePlaceNotInitialized
	}

	if orderID == uuid.Nil {
		return errs.NewValueIsInvalidError("orderID")
	}

	if volume <= 0 {
		return errs.NewValueIsInvalidErrorWithCause("volume", ErrVolumeMustPositive)
	}

	if sp.IsOccupied() {
		return ErrStoragePlaceIsOccupied
	}

	if sp.totalVolume < volume {
		return errs.NewValueIsOutOfRangeError("volume", volume, 1, sp.totalVolume)
	}

	sp.orderID = &orderID
	return nil
}

func (sp *StoragePlace) Clear(orderID uuid.UUID) error {
	if sp == nil {
		return ErrStoragePlaceNotInitialized
	}

	if sp.orderID == nil || *sp.orderID != orderID {
		return errs.NewObjectNotFoundError("orderID", orderID)
	}

	sp.orderID = nil
	return nil
}

func (sp *StoragePlace) IsOccupied() bool {
	if sp == nil {
		return false
	}

	return sp.orderID != nil
}
