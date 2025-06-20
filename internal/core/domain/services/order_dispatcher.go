package services

import (
	"errors"
	"math"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
)

type OrderDispatcher interface {
	Dispatch(*order.Order, []*courier.Courier) (*courier.Courier, error)
}

var (
	ErrCantAssignOrder = errors.New("cannot assign the order")
	ErrNoRightCourier  = errors.New("no right courier")
)

var _ OrderDispatcher = (*orderDispatcher)(nil)

type orderDispatcher struct{}

func NewOrderDispatcher() OrderDispatcher { return new(orderDispatcher) }

func (o *orderDispatcher) Dispatch(ordering *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if ordering == nil {
		return nil, errors.Join(ErrCantAssignOrder, errs.NewValueIsRequiredError("order"))
	}

	if ordering.Status() != order.StatusCreated {
		return nil, errors.Join(ErrCantAssignOrder, errs.NewExpectationFailedError("ordering.status", ordering.Status(), order.StatusCreated))
	}

	index, remain := -1, math.MaxFloat64
	for i := range couriers {
		ok, err := couriers[i].CanTakeOrder(ordering)
		if ok && err == nil {
			dt, err := couriers[i].CalculateTimeToLocation(ordering.Location())
			if dt < remain && err == nil {
				index, remain = i, dt
			}
		}
	}

	if index >= 0 {
		err := couriers[index].TakeOrder(ordering)
		if err != nil {
			return nil, errors.Join(ErrCantAssignOrder, err)
		}

		return couriers[index], nil
	}

	return nil, errors.Join(ErrCantAssignOrder, ErrNoRightCourier)
}
